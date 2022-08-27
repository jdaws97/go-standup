package logic

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jdaws97/go-standup/config"
)

var DATE string = time.Now().Format("2006-01-02")

type Standup struct {
	Category string
	Config config.Config
	Days_ago int
	Sentence []string
}


func parse_error(err error) {
	if err != nil {
		log.Fatal(err)
	} else {
		return
	}

}


func isElementExist(s []string, str string) bool {
	for _, v := range s {
	  if v == str {
		return true
	  }
	}
	return false
  }


func Check_path(path string, config_struct config.Config) {
	file_info, err := os.Stat(path)
	parse_error(err)
	config_info, err := os.Stat(config_struct.File_path)
	parse_error(err)

	if file_info.IsDir() {
		err := os.Chdir(path)
		parse_error(err)
	} else {
		if config_info.IsDir() {
			err := os.Chdir(config_struct.File_path)
			parse_error(err)
		}
	}

}


func Check_standup(config_struct config.Config, standup Standup) {

	Check_path(standup.Config.File_path, standup.Config)
	var file_name string

	files, err := ioutil.ReadDir(".")
	parse_error(err)
	for _, f := range files {
		file := f.Name()
		if standup.Days_ago == 0 {
			if file == fmt.Sprintf("standup_%s.txt", DATE) {
				file_name = file
			}
		} else {
			day := standup.Days_ago - (standup.Days_ago * 2)
			new_date := time.Now().AddDate(0, 0, day).Format("2006-01-02")
			if file == fmt.Sprintf("standup_%s.txt", new_date) {
				file_name = file
			}
		}
	}

	if strings.TrimSpace(file_name) != "" {
		return 
	} else if strings.TrimSpace(file_name) == "" && standup.Days_ago == 0 {
		Create_standup(standup.Config)
	} else {
		log.Fatalf("There isn't a file from %v days ago!", standup.Days_ago)
	}


}


func Append_standup(config_struct config.Config, standup Standup) {

	Check_path(standup.Config.File_path, standup.Config)
	var lines []string

	var sentence string = strings.Join(standup.Sentence, "")


	if standup.Days_ago == 0 {
		f, err := os.Open(fmt.Sprintf("standup_%s.txt", DATE))
		parse_error(err)
		scanner := bufio.NewScanner(f)

		start := -1
		var category_to_stop string
		for scanner.Scan() {
			line := scanner.Text()
			formatted_line := fmt.Sprintf("%s\n", line)
			lines = append(lines, formatted_line)
		}
		for _, current_line := range lines {
			if strings.Contains(current_line, standup.Category) {
				categories := standup.Config.Categories
				for index, _ := range categories {
					if strings.Contains(categories[index], standup.Category) {
						if standup.Category == "NOTES" {
							category_to_stop = "stop"
						} else {
							category_to_stop = categories[index+1]
						}
					}
				}
				start += 1
				break
			}
			start += 1
		}

		var current_line string
		i := -1
		new_lines := lines[start:]
		for _, line := range new_lines {
			current_line = strings.TrimSpace(line)
			if strings.Contains(current_line, category_to_stop) {
				break
			}
			i += 1
		}
		place_line := i + start
		lines[place_line] = fmt.Sprintf("\t[%v]: %s\n\n", i-1, sentence)

		err = ioutil.WriteFile(f.Name(), []byte(strings.Join(lines, "")), 0644)
		parse_error(err)
		fmt.Printf("Added your sentence to category %s! Here is how your current notes look:\n\n", standup.Category)
		for _, current_line := range lines {
			fmt.Printf("%s", current_line)
		}

	} else {
		log.Fatal("You can't update an older file!")
	}

}


func Open_standup(config_struct config.Config, standup Standup) {

	if standup.Days_ago == 0 {
		cmd := exec.Command("vim", fmt.Sprintf("%s/standup_%s.txt", config_struct.File_path, DATE))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		parse_error(err)
	} else {
		day := standup.Days_ago - (standup.Days_ago * 2)
		new_date := time.Now().AddDate(0, 0, day).Format("2006-01-02")
		cmd := exec.Command("vim", fmt.Sprintf("%s/standup_%s.txt", config_struct.File_path, new_date))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		parse_error(err)
	}
}


func Create_standup(config_struct config.Config) {

	f, err := os.Create(fmt.Sprintf("standup_%s.txt", DATE))
	parse_error(err)
	var categories_to_write []string
	for _, category :=  range config_struct.Categories {
		categories_to_write = append(categories_to_write, fmt.Sprintf("%s\n\n\n", category))
	}
	err = ioutil.WriteFile(f.Name(), []byte(strings.Join(categories_to_write, "")), 0644)
	parse_error(err)

}


func Remove_old_standups(config_struct config.Config) {

	var dates []string
	for i := 0; i <=config_struct.Days-1; i++ {
		new_time := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates = append(dates, new_time)
	}
	Check_path(config_struct.File_path, config_struct)

	var file_list []string
	files, err := ioutil.ReadDir(".")
	parse_error(err)
	for _, f := range files {
		file := f.Name()
		if strings.Contains(file, "standup") && file != ".standup-config.json" {
			for _, date := range dates {
				if strings.Contains(file, date) {
					file_list = append(file_list, file)
					continue
				}
			}
		}
	}

	for _, f := range files {
		file := f.Name()
		if strings.Contains(file, "standup") && file != ".standup-config.json" {
			if !isElementExist(file_list, file) {
				os.Remove(file)
			}
		}
	}

}
