package main

import (
	"flag"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"os"
	"strings"
	"sync"
	"time"
	jsonlog "todo_list_be/internal/log"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
	security struct {
		jwtSecret string
		exp       time.Duration
	}
}
type application struct {
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
	db     *gorm.DB
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "server port")
	flag.StringVar(&cfg.db.dsn, "db-dns", "root:123456@tcp(127.0.0.1:3309)/todo_list?charset=utf8mb4&parseTime=True&loc=Local", "Mysql Dsn")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "MySql max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "MySql max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle_time", "15m", "MySql max connection idle time")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.StringVar(&cfg.security.jwtSecret, "security-jwt-secret", "dsfdsfdsfdsf", "jwt secret key")
	flag.DurationVar(&cfg.security.exp, "security-exp", time.Minute*30, "jwt exp")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := gorm.Open(mysql.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}
	err = migrate(db)

	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		db:     db,
	}
	//test, _ := hashPassword("123456")
	//fmt.Println(test)
	err = app.serve()

	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
