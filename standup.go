package main

import (
	"log"
	"os"
	"strings"

	config "github.com/jdaws97/go-standup/config"
	logic "github.com/jdaws97/go-standup/logic"
	"github.com/jessevdk/go-flags"
)

var HOME = os.Getenv("HOME")
var initial_config = &config.Config{7, []string{"DONE", "IN-PROGRESS", "BLOCKERS", "NOTES"}, HOME}


var opts struct {
	Category string `short:"c" long:"category" description:"Category you'd like your notes to be placed under. See your current categories within your config!\nEXAMPLE: standup -c in-progress"`

	Config bool `long:"config" description:"Open your config"`

	Open bool `long:"open" description:"Open your standup notes"`

	Sentence []string `short:"s" description:"The sentence you want added to your notes"`

	DaysAgo int `short:"d" default:"0" description:"Number of days ago for checking old standup notes\n Here's an example command:\n standup -c in-progress -s Finished ticket 1234"`
}


func parse_error(err error) {
	if err != nil {
		log.Fatal(err)
	} else {
		return
	}

}


func map_category(category string, category_config []string) string {
	
	for _, categ := range category_config {
		category = strings.ToUpper(category)
		if categ == category {
			return category
		}
	}
	return ""
}


func main() {

	// Initialize everything! 
	parsed_config := config.Check_config(initial_config)
	var category string

	// Initialize parsed arguments into variables
	_, err := flags.Parse(&opts)
	parse_error(err)

	// Determine allowable categories based on config
	if strings.TrimSpace(opts.Category) != "" {
		category = map_category(opts.Category, parsed_config.Categories)
		if strings.TrimSpace(category) == "" {
			log.Fatal("The category fsrom your argument is not in your config!")
		}
	} else {
		category = "NOTES"
	}

	// Open the config if the open command isn't given
	if opts.Config && !opts.Open {
		config.Open_config()
		parsed_config = config.Check_config(initial_config)
	}

	// Warn the user you can't open your notes and the config at the same time
	if opts.Open && opts.Config {
		log.Print("You can't open your notes and the config at the same time, choose one or the other!")
	}

	// Create the standup struct holding all the values needed to be passed down for the notes
	standup := logic.Standup{category, parsed_config, opts.DaysAgo, opts.Sentence}
	logic.Check_standup(parsed_config, standup)

	// If there's a sentence argument, append it to the standup notes
	if len(opts.Sentence) != 0 {
		if opts.DaysAgo > 0 {
			log.Fatalf("You can't modify a standup from %v days ago!", opts.DaysAgo)
		}
		logic.Append_standup(parsed_config, standup)
	// Open the standup notes if open argument given and not config
	} else if opts.Open && !opts.Config {
		logic.Open_standup(parsed_config, standup)
	} 

	// Remove old standups based on the number of days in the config!
	if !opts.Config {
		logic.Remove_old_standups(standup.Config)
	}
}