package main

import (
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/jghiloni/do-all/repeater"
)

var (
	// Version the current version of the plugin
	Version = "1.0.0"
)

func main() {
	doAll := &repeater.Repeater{
		Version: Version,
		Writer:  os.Stdout,
	}

	plugin.Start(doAll)
}
