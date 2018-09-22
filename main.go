package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bojand/ghz-web/config"
)

var (
	// set by goreleaser with -ldflags="-X main.version=..."
	version = "dev"
	cPath   = flag.String("config", "", "Path to the config file.")
	v       = flag.Bool("v", false, "Print the version.")
)

var usage = `Usage: ghz-web [options...]
Options:
  -config	Path to the config JSON file.
  -v  Print the version.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}

	flag.Parse()

	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	cfgPath := strings.TrimSpace(*cPath)

	conf, err := config.Read(cfgPath)
	if err != nil {
		panic(err)
	}

	info := &config.Info{
		Version:   version,
		GOVersion: runtime.Version(),
		StartTime: time.Now(),
	}

	app := Application{
		Config: conf,
		Info:   info,
	}

	app.Start()
}
