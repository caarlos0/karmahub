package main

import (
	"github.com/caarlos0/karmahub/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
