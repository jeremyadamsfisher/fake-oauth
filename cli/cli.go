package cli

import (
	"fmt"
	"os"

	"github.com/fred1268/go-clap/clap"
)

type Config struct {
	DownstreamURL string `clap:"--downstream,mandatory"`
}

func CommandLineInterface() *Config {
	cfg := &Config{}
	if _, err := clap.Parse(os.Args[1:], cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}
