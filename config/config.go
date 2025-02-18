package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const conf = "conf.json"

type Json struct {
	LastTemplate string `json:"last_template"`
	LastGenerate string `json:"last_generate"`
}

func (j *Json) Load() {
	data, err := os.ReadFile(Path())
	if err == nil {
		err = json.Unmarshal(data, j)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func (j *Json) Update() {
	data, err := json.Marshal(j)
	if err == nil {
		if err = os.WriteFile(Path(), data, os.ModePerm); err == nil {
			return
		}
	}
	log.Fatalln("config.Update()", err)
}

func ConfigDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalln("Error getting executable path:", err)
	}
	executableName := filepath.Base(executablePath)
	return filepath.Join(".config", executableName)
}

func Path() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ConfigDir(), conf)
	} else {
		log.Fatalln("config.Path()", err)
	}
	return ""
}
