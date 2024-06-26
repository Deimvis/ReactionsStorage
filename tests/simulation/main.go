package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/Deimvis/reactionsstorage/tests/simulation/src/configs"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/metrics"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/models"
	rs "github.com/Deimvis/reactionsstorage/tests/simulation/src/rs_client"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/utils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to simulation config (JSON or YAML)")
	quiet := flag.Bool("q", false, "Disable logging")
	flag.Parse()

	logger := newLogger(*quiet)
	zap.ReplaceGlobals(logger.Desugar())
	defer logger.Sync()

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	logger.Infof("Config:\n%s", string(configData))

	config := configs.Simulation{}
	decoder := yaml.NewDecoder(bytes.NewReader(configData))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}
	logger.Debugf("Parsed config:\n%s", spew.Sdump(config))

	rand.Seed(config.Seed)

	var pgwRecorder metrics.HTTPRecorder = nil
	if config.PrometheusPushgateway != nil {
		pgwRecorder = metrics.NewPrometheusPushgatewayRecorder(
			config.PrometheusPushgateway.Host,
			config.PrometheusPushgateway.Port,
			config.PrometheusPushgateway.SSL,
			"rs_client",
		)
		defer pgwRecorder.Sync()
		cron := utils.NewCron(func() {
			pgwRecorder.Sync()
		}, time.Duration(config.PrometheusPushgateway.PushIntervalS)*time.Second)
		cron.Start()
		defer cron.Stop()
	}
	rsClient := rs.NewClientHTTP(config.Server.Host, config.Server.Port, config.Server.SSL, logger, pgwRecorder)

	var topics []models.Topic
	for i, topicConf := range config.Rules.Topics {
		namespace := models.NewNamespace(topicConf.NamespaceId, rsClient)
		for j := 0; j < topicConf.Count; j++ {
			id := fmt.Sprintf("topic_%d_%d", i, j)
			topics = append(topics, models.NewTopic(id, namespace, topicConf.Size, topicConf.ShufflePerUser))
		}
	}

	var apps []models.App
	for i := 0; i < config.Rules.Users.Count; i++ {
		app := models.NewApp(rsClient, topics, config.Rules.Users.Screen.VisibleEntitiesCount, logger)
		apps = append(apps, app)
	}

	setupSigHandlers()
	run(config, apps, logger)
}

func run(config configs.Simulation, apps []models.App, logger *zap.SugaredLogger) {
	n := len(apps)

	wgChs := make([]chan *sync.WaitGroup, n)
	userTurnInds := make([]atomic.Uint64, n)
	for i, a := range apps {
		wgChs[i] = make(chan *sync.WaitGroup)
		go runUser(config, a, wgChs[i], &userTurnInds[i], logger)
	}

	expectedTurnDur := time.Duration(config.Rules.Turns.MinDurMs) * time.Millisecond
	var wg *sync.WaitGroup
	for i := 0; i < config.Rules.Turns.Count; i++ {
		logger.Infof("Turn %d", i)

		timer := time.NewTimer(expectedTurnDur)
		wg = &sync.WaitGroup{}
		wg.Add(n)
		for _, ch := range wgChs {
			ch <- wg
		}
		<-timer.C

		logUsersBehind(logger, i, userTurnInds, 3)
	}
	extraSimDur := utils.MeasureDuration(func() {
		wg.Wait()
	})
	expectedSimDur := expectedTurnDur * time.Duration(config.Rules.Turns.Count)
	if float64(extraSimDur) > float64(0.05)*float64(expectedSimDur) {
		logger.Warnf("Simulation was significantly longer than expected: %s (expected: %s)", expectedSimDur+extraSimDur, expectedSimDur)
	}

	logger.Info("Simulation finished")
}

func logUsersBehind(logger *zap.SugaredLogger, turn int, userTurnInds []atomic.Uint64, maxTurnsBehind int) {
	usersBehind := 0
	for j := range userTurnInds {
		if turn >= maxTurnsBehind && uint64(turn-maxTurnsBehind) > userTurnInds[j].Load() {
			usersBehind++
		}
	}
	if usersBehind > 0 {
		logger.Warnf("%d users are more than %d turns behind", usersBehind, maxTurnsBehind)
	}
}

func runUser(config configs.Simulation, app models.App, wgCh <-chan *sync.WaitGroup, turnInd *atomic.Uint64, logger *zap.SugaredLogger) {
	// simulates app loading
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	user := logIn(config, app, logger)

	for i := 0; i < config.Rules.Turns.Count; i++ {
		wg := <-wgCh
		turnInd.Store(uint64(i))
		time.Sleep(time.Duration(rand.Intn(config.Rules.Users.TurnStartSkewMs)) * time.Millisecond)
		if needRefresh(config, i) {
			// do not wait, it's a background process
			app.Refresh(user.GetId())
		}

		start := time.Now()
		{
			id := user.GetId() // it's unsafe to use user after quit
			user.DoRandomAction()
			logger.Infof("User %s finished turn %d (%s)", id, i, time.Since(start))
		}

		if user.IsQuit() {
			// user quit => log in a new user
			user = logIn(config, app, logger)
		}
		wg.Done()
	}
}

// logIn logins a new user
func logIn(config configs.Simulation, app models.App, logger *zap.SugaredLogger) models.User {
	id := createUserId(config, userCounter.Add(1))
	user := models.NewUser(id, app, config.Rules.Users.ActionProbs, logger)
	app.Restart(id).Wait()
	logger.Infof("User %s logged in", id)
	return user
}

// needRefresh returns whether it's right turn to call Refresh
func needRefresh(config configs.Simulation, turn int) bool {
	return turn%config.Rules.Users.App.Background.RefreshReactions.TimerInTurns == 0
}

func createUserId(config configs.Simulation, ind uint64) string {
	template := defaultUserIdTemplate
	if config.Rules.Users.IdTemplate != nil {
		template = *config.Rules.Users.IdTemplate
	}
	return fmt.Sprintf(template, ind)
}

func newLogger(quiet bool) *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.Sampling = nil
	level := zap.InfoLevel
	if utils.IsDebugEnv() {
		level = zap.DebugLevel
	}
	if quiet {
		level = zap.FatalLevel
	}
	config.Level.SetLevel(level)
	return zap.Must(config.Build()).Sugar()
}

func setupSigHandlers() {
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			debug.PrintStack()
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}

var userCounter atomic.Uint64
var defaultUserIdTemplate = "user_%06d"
