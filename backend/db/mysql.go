package db

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/kylelemons/go-gypsy/yaml"
)

const (
	dbSessionTimeout = 10 * time.Second
	dbMaxConnections = 50
)

var (
	mutex    = sync.RWMutex{}
	sessions = make(map[string]*dbr.Session)
)

func GetSession(hint string) *dbr.Session {
	mutex.RLock()
	session, ok := sessions[hint]
	mutex.RUnlock()
	if ok {
		return session
	}

	mutex.Lock()
	defer mutex.Unlock()
	sessions[hint] = newSession()
	return sessions[hint]
}

func newSession() *dbr.Session {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	// dsn := os.Getenv("MYSQL_DSN")
	dsl := getProperties()
	if dsl == "" {
		panic("MYSQL_DSL environment variable is not set")
	}

	conn, err := dbr.Open("mysql", dsl, nil)
	if err != nil {
		panic(err)
	}
	s := conn.NewSession(nil)
	s.Timeout = dbSessionTimeout
	s.SetMaxOpenConns(dbMaxConnections)
	s.SetMaxIdleConns(dbMaxConnections)
	return s
}

// dslを返却
func getProperties() string {
	cfgFile := filepath.Join("./db/dbconfig.yml")

	f, err := yaml.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}

	open, err := f.Get("development.datasource")
	if err != nil {
		panic(err)
	}
	return os.ExpandEnv(open)
}
