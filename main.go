package main

import (
	cmd "github.com/jdamata/bb-bulk-ssh-key/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
