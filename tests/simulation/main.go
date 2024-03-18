package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/Deimvis/reactionsstorage/tests/simulation/configs"
	"github.com/Deimvis/reactionsstorage/tests/simulation/models"
	rs "github.com/Deimvis/reactionsstorage/tests/simulation/rs_client"
	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

func main() {
	logger := newLogger()
	zap.ReplaceGlobals(logger.Desugar())
	defer logger.Sync()
	configPath := flag.String("config", "config.yaml", "Path to simulation config (JSON or YAML)")
	flag.Parse()

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	logger.Infof("Config:\n%s", string(configData))

	config := configs.Simulation{}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}
	logger.Debugf("Parsed config:\n%s", spew.Sdump(config))

	rand.Seed(config.Seed)

	rsClient := rs.NewClientHTTP(config.Server.Host, config.Server.Port, config.Server.SSL)

	var topics []*models.Topic
	for i, topicConf := range config.Rules.Topics {
		namespace := models.NewNamespace(topicConf.NamespaceId, rsClient)
		for j := 0; j < topicConf.Count; j++ {
			id := fmt.Sprintf("topic_%d_%d", i, j)
			topics = append(topics, models.NewTopic(id, topicConf.Size, namespace))
		}
	}

	var users []models.User
	for i, userConf := range config.Rules.Users {
		app := models.NewApp(rsClient, topics, userConf.Screen.VisibleEntitiesCount)
		id := fmt.Sprintf("user_%d", i)
		users = append(users, models.NewUser(id, app))
	}

	run(config, users, logger)
}

func run(config configs.Simulation, users []models.User, logger *zap.SugaredLogger) {
	n := len(users)

	wgChs := make([]chan *sync.WaitGroup, n)
	for i, u := range users {
		wgChs[i] = make(chan *sync.WaitGroup)
		go runUser(u, config.Rules.TurnCount, wgChs[i])
	}

	for i := 0; i < config.Rules.TurnCount; i++ {
		logger.Infof("Start turn %d\n", i)
		timer := time.NewTimer(time.Duration(config.Rules.TurnDurS) * time.Second)
		wg := &sync.WaitGroup{}
		wg.Add(n)
		for _, ch := range wgChs {
			ch <- wg
		}
		<-timer.C
		// TODO: record wait time and log warning if it's significant
		wg.Wait()
	}

	fmt.Println("simulation finished")
}

func runUser(user models.User, turnCount int, wgCh <-chan *sync.WaitGroup) {
	for i := 0; i < turnCount; i++ {
		wg := <-wgCh
		fmt.Printf("user %s finished turn %d\n", user.GetId(), i)
		wg.Done()
	}
}

func newLogger() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.Sampling = nil
	level := zap.InfoLevel
	if utils.IsDebugEnv() {
		level = zap.DebugLevel
	}
	config.Level.SetLevel(level)
	return zap.Must(config.Build()).Sugar()
}
