package cmd

import (
	"bufio"
	"errors"
	"github.com/leigme/gft/config"
	"github.com/leigme/gft/model"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	p  = model.Param{}
	cj = config.Json{}
)

var rootCmd = &cobra.Command{
	Use:   "gft",
	Short: "Generate files through templates",
	Long:  `Generate files through templates, for example: --t your template --g create file --a template param`,
	Run: func(cmd *cobra.Command, args []string) {
		defer cj.Update()
		bindLast()
		if err := paramCheck(); err != nil {
			log.Fatalln(err)
		}
		generate()
		cj.LastTemplate = p.Template
		cj.LastGenerate = p.Generate
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	createDir()
	cj.Load()
	rootCmd.PersistentFlags().StringVar(&p.Template, "t", cj.LastTemplate, "")
	rootCmd.PersistentFlags().StringVar(&p.Args, "a", "", "")
	rootCmd.PersistentFlags().StringVar(&p.Generate, "g", cj.LastGenerate, "")
}

func createDir() {
	configPath := config.Path()
	_, err := os.Stat(filepath.Dir(configPath))
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
			if err == nil {
				return
			}
		}
		log.Fatalln(err)
	}
}

func bindLast() {
	if strings.EqualFold(p.Template, "") && !strings.EqualFold(cj.LastTemplate, "") {
		p.Template = cj.LastTemplate
	}
	if strings.EqualFold(p.Generate, "") && !strings.EqualFold(cj.LastGenerate, "") {
		p.Generate = cj.LastGenerate
	}
}

func paramCheck() error {
	if strings.EqualFold(p.Template, "") {
		return errors.New("--t is nil")
	}
	if strings.EqualFold(p.Generate, "") {
		return errors.New("--g is nil")
	}

	if strings.EqualFold(p.Args, "") {
		return errors.New("--a is nil")
	}
	if !strings.Contains(p.Args, ":") {
		return errors.New("--a must be contains `:`")
	}
	return nil
}

func generate() {
	data, err := os.ReadFile(p.Template)
	if err != nil {
		log.Fatalln("path: ", p.Template, "read file fail", err)
	}
	t, err := template.New("gpf").Parse(string(data))
	if err != nil {
		log.Fatalln("path: ", p.Template, "parse template fail", err)
	}
	if _, err = os.Stat(p.Generate); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln("path: ", p.Generate, "generate file fail", err)
		}
		if err = os.MkdirAll(filepath.Dir(p.Generate), os.ModePerm); err != nil {
			log.Println(err)
		}
	}
	f, err := os.OpenFile(p.Generate, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	w := bufio.NewWriter(f)
	if err = t.Execute(w, paramMap(p.Args)); err != nil {
		log.Fatalln(err)
	}
	if err = w.Flush(); err != nil {
		log.Fatalln(err)
	}
}

func paramMap(arg string) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	args := strings.Split(arg, ",")
	if len(args) > 0 {
		for _, s := range args {
			v := strings.Split(s, ":")
			if len(v) == 2 {
				result[v[0]] = v[1]
			}
		}
	}
	return result
}
