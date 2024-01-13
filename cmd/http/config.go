package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT,5000"`
	DBPassword string `env:"DB_PASSWORD,required"`
}

func (c Config) SPrint() (envs string) {
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := strings.Split(field.Tag.Get("env"), ",")
		name := envTag[0]
		value := envTag[1]
		envs += fmt.Sprintf("%s - %s\n", name, value)
	}
	return
}

func (c Config) loadFromEnv() (conf Config) {
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := strings.Split(field.Tag.Get("env"), ",")
		envName := envTag[0]
		defaultValue := envTag[1]
		value := os.Getenv(envName)
		if value == "" && defaultValue != "required" {
			f := reflect.ValueOf(&conf).Elem().FieldByName(field.Name)
			f.SetString(defaultValue)
		} else {
			f := reflect.ValueOf(&conf).Elem().FieldByName(field.Name)
			f.SetString(value)
		}
	}
	return
}

func (c Config) validate() {
	var validationMsg string
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		envTag := strings.Split(t.Field(i).Tag.Get("env"), ",")
		envName := envTag[0]
		envValue := envTag[1]
		if envValue == "required" && value.String() == "" {
			validationMsg += fmt.Sprintf("%s is required\n", envName)
		}
	}
	if len(validationMsg) != 0 {
		panic(validationMsg)
	}
}

func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config := Config{}
	config = config.loadFromEnv()
	config.validate()
	return config
}
