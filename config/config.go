package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var HOME = os.Getenv("HOME")

type Config struct{
	Days int
	Categories []string
	File_path string
}


func parse_error(err error) {
	if err != nil {
		log.Fatal(err)
	} else {
		return
	}
}


func Open_config() {
	cmd := exec.Command("vim", fmt.Sprintf("%s/.standup-config.json", HOME))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	parse_error(err)
}


func walk_directory(file_name string, path string) (bool) {
	files, err := ioutil.ReadDir(path)
	parse_error(err)
	for _, f := range files {
		file := f.Name()
		if strings.Contains(file, file_name) {
			return true
		}
	}
	return false
}


func Check_config(initial_config *Config) (Config){

	file_bool := walk_directory(".standup-config.json", HOME)

	if file_bool {
		return parse_config()
	} else {
		create_config(initial_config)
	}
	return parse_config()
}


func create_config(initial_config *Config) {

	os.Chdir(HOME)

	created_file, err := os.Create(".standup-config.json")
	parse_error(err)

	result, err := json.MarshalIndent(initial_config, "", "  ")
	parse_error(err)

	err = ioutil.WriteFile(created_file.Name(), result, 0644)
	parse_error(err)

	log.Printf("Config %v created at path %v", created_file.Name(), HOME)
}


func parse_config() (Config){

	os.Chdir(HOME)

	config_content, err := ioutil.ReadFile("./.standup-config.json")
	parse_error(err)

	var payload Config

	err = json.Unmarshal(config_content, &payload)
	parse_error(err)

	return payload
}

