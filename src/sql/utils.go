package sql

import (
	"os"
	"path/filepath"
)

func ReadQueryFile(queryRelPath string) string {
	path := filepath.Join(os.Getenv("SQL_SCRIPTS_DIR"), queryRelPath)
	query, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(query)
}
