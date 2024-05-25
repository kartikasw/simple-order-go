package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      App
	Database Database
}

func NewConfig(v *viper.Viper) Config {
	return Config{
		App:      NewApp(v),
		Database: NewDatabase(v),
	}
}

type App struct {
	Port int
	Host string
}

func NewApp(v *viper.Viper) App {
	return App{
		Port: v.GetInt("app.port"),
	}
}

type Database struct {
	Name     string
	Host     string
	Port     int
	Password string
	User     string
	Timezone string
	SslMode  string
}

func NewDatabase(v *viper.Viper) Database {
	return Database{
		Name:     v.GetString("database.name"),
		Host:     v.GetString("database.host"),
		Port:     v.GetInt("database.port"),
		Password: v.GetString("database.password"),
		User:     v.GetString("database.user"),
		Timezone: v.GetString("database.timezone"),
		SslMode:  v.GetString("database.sslmode"),
	}
}

func LoadConfig(path string) Config {
	v := viper.New()
	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Load config error: ", err.Error())
	}

	return NewConfig(v)
}
