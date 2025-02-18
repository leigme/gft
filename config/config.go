package config

import (
	"bufio"
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
}

func (j *Json) Update() {
	data, err := json.Marshal(j)
	if err == nil {
		f, err := os.OpenFile(Path(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err == nil {
			defer func() {
				err = f.Close()
				if err != nil {
					log.Fatalln("config.Update()", err)
				}
			}()
			wr := bufio.NewWriter(f)
			_, err = wr.Write(data)
			err = wr.Flush()
			return
		}
	}
}

func Dir() string {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalln("Error getting executable path:", err)
	}
	executableName := filepath.Base(executablePath)
	return filepath.Join(".config", executableName)
}

func Path() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, Dir(), conf)
	} else {
		log.Fatalln("config.Path()", err)
	}
	return ""
}
