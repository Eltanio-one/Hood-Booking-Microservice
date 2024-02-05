package config

import (
	"encoding/json"
	"os"
)

// define the config struct, which has integrated structs to be used in JSON format when read.
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
}

// ReadConfigFile takes a filename as a string and returns a Config struct object and an error.
// This function is used to read the filename given, it is expected that the filename will be the same as the json file within the config package.
func ReadConfigFile(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return config, err
	}
	return config, err
}
