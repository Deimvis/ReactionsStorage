package sql

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

func ReadQueryFile(queryRelPath string) string {
	path := filepath.Join(os.Getenv("SQL_SCRIPTS_DIR"), queryRelPath)
	query, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(query)
}

// ParseQueries parses sql queries separated with ";"
func ParseQueries(s string) []string {
	queries := strings.Split(s, ";")
	utils.FilterIn(&queries, func(q string) bool { return strings.TrimSpace(q) != "" })
	return queries
}
