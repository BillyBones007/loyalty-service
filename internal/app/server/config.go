package server

import (
	"errors"
	"flag"
	"log"

	"github.com/caarlos0/env"
)

// Custom errors
var (
	errSecretKey   error = errors.New("env var SECRET_KEY is empty")
	errAddrDB      error = errors.New("dsn string connection to the database not found")
	errAddrAccrual error = errors.New("accrual system address not found")
)

// Main configuration structure
type Config struct {
	SecretKey   string `env:"SECRET_KEY"`
	AddrServ    string `env:"RUN_ADDRESS"`
	AddrDB      string `env:"DATABASE_URI"`
	AddrAccrual string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// Commandline flags
type cmdFlags struct {
	runAddr     string
	dbAddr      string
	accrualAddr string
}

// Initialization function
func initConfig() *Config {
	conf := Config{}
	err := env.Parse(&conf)
	if err != nil {
		log.Fatal(err)
	}
	flags := cmdFlags{}
	flag.StringVar(&flags.runAddr, "a", ":8081", "Run address server. By default localhost:8081.")
	flag.StringVar(&flags.dbAddr, "d", "", "DSN string connection to the database.")
	flag.StringVar(&flags.accrualAddr, "r", "", "Accrual system addreess.")
	flag.Parse()

	if conf.SecretKey == "" {
		// At the time of development, the secret key is specified manually
		// log.Fatal(errSecretKey)
		conf.SecretKey = "this is secret key"
	}
	if conf.AddrDB == "" {
		if flags.dbAddr == "" {
			log.Fatal(errAddrDB)
		} else {
			conf.AddrDB = flags.dbAddr
		}
	}
	if conf.AddrServ == "" {
		conf.AddrServ = flags.runAddr
	}
	if conf.AddrAccrual == "" {
		if flags.accrualAddr == "" {
			log.Fatal(errAddrAccrual)
		} else {
			conf.AddrAccrual = flags.accrualAddr
		}
	}
	return &conf
}
