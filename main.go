/*
Copyright © 2025 leig <leigme@gmail.com>
*/
package main

import (
	"log"

	"github.com/leigme/gft/cmd"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	cmd.Execute()
}
