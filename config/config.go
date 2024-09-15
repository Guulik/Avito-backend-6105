package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	SERVER_ADDRESS    string `yaml:"address"`
	POSTGRES_CONN     string
	POSTGRES_JDBC_URL string
	POSTGRES_USERNAME string `yaml:"user"`
	POSTGRES_PASSWORD string `yaml:"password"`
	POSTGRES_HOST     string `yaml:"host"`
	POSTGRES_PORT     string `yaml:"port"`
	POSTGRES_DATABASE string `yaml:"DBName"`

	ENV string
}

func MustLoad() *Config {
	config := &Config{}

	if v := os.Getenv("DEV_ENV"); v != "" {
		config.ENV = v
	} else {
		config.ENV = "prod"
	}
	fmt.Println(config.ENV)

	//if local or stage, read from file
	if config.ENV == "local" {
		config = MustLoadPath("./config/local.yaml")
	} else if config.ENV == "stage" {
		config = MustLoadPath("./config/stage.yaml")
	} else {
		//else read from env
		config = MustGetEnv()
	}

	fmt.Println(config)
	return config
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("cannot read config: " + err.Error())
	}

	//compile connection string
	postgresConn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.POSTGRES_USERNAME, config.POSTGRES_PASSWORD, config.POSTGRES_HOST, config.POSTGRES_PORT, config.POSTGRES_DATABASE)

	//compile jdbc url
	postgresJDBCUrl := fmt.Sprintf(
		"jdbc:postgresql://%s:%s/%s",
		config.POSTGRES_HOST, config.POSTGRES_PORT, config.POSTGRES_DATABASE)

	//set config conn string
	config.POSTGRES_CONN = postgresConn
	//set config jdbc url
	config.POSTGRES_JDBC_URL = postgresJDBCUrl

	//one more time, because cleanenv.ReadConfig rewrites it to ""
	config.ENV = "local"
	return &config
}

func MustGetEnv() *Config {
	config := &Config{}

	envs := map[string]*string{
		"SERVER_ADDRESS":    &config.SERVER_ADDRESS,
		"POSTGRES_CONN":     &config.POSTGRES_CONN,
		"POSTGRES_JDBC_URL": &config.POSTGRES_JDBC_URL,
		"POSTGRES_USERNAME": &config.POSTGRES_USERNAME,
		"POSTGRES_PASSWORD": &config.POSTGRES_PASSWORD,
		"POSTGRES_HOST":     &config.POSTGRES_HOST,
		"POSTGRES_PORT":     &config.POSTGRES_PORT,
		"POSTGRES_DATABASE": &config.POSTGRES_DATABASE,
	}

	for env, ptr := range envs {
		if v := os.Getenv(env); v != "" {
			*ptr = v
		} else {
			log.Fatal(fmt.Sprintf("failed to get env: %s", env))
		}
	}

	config.ENV = "prod"
	return config
}
