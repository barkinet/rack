package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/util"
	"gopkg.in/ini.v1"
)

func configure(c *cli.Context) {
	reader := bufio.NewReader(os.Stdin)
	m := map[string]string{
		"username": "",
		"apikey":   "",
		"region":   "",
	}
	fmt.Print("Rackspace Username: ")
	username, _ := reader.ReadString('\n')
	m["username"] = strings.TrimSuffix(username, string('\n'))

	fmt.Print("Rackspace API key: ")
	apiKey, _ := reader.ReadString('\n')
	m["apikey"] = strings.TrimSuffix(apiKey, string('\n'))

	fmt.Print("Rackspace Region: ")
	region, _ := reader.ReadString('\n')
	m["region"] = strings.TrimSuffix(region, string('\n'))

	fmt.Print("Profile Name: ")
	profile, _ := reader.ReadString('\n')
	profile = strings.TrimSuffix(profile, string('\n'))

	configFile, err := configFile()
	var cfg *ini.File
	cfg, err = ini.Load(configFile)
	if err != nil {
		fmt.Printf("Error loading config file: %s\n", err)
		cfg = ini.Empty()
	}

	if profile == "" {
		profile = "DEFAULT"
	}
	section, err := cfg.NewSection(profile)
	if err != nil {
		fmt.Printf("Error creating new section [%s] in config file: %s\n", profile, err)
		return
	}

	for key, val := range m {
		section.NewKey(key, val)
	}

	err = cfg.SaveTo(configFile)
	if err != nil {
		fmt.Printf("Error saving config file: %s\n", err)
		return
	}
}

func configFile() (string, error) {
	dir, err := util.RackDir()
	if err != nil {
		return "", fmt.Errorf("Error reading from cache: %s", err)
	}
	filepath := path.Join(dir, "config")
	// check if the cache file exists
	if _, err := os.Stat(filepath); err == nil {
		return filepath, nil
	}
	// create the cache file if it doesn't already exist
	f, err := os.Create(filepath)
	defer f.Close()
	return filepath, err
}
