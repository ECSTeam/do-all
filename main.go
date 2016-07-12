package main

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/do-all/repeater"
)

var (
	// Version the current version of the plugin
	Version = "1.0.0"
)

func main() {
	doAll := &repeater.Repeater{
		Version: Version,
	}

	plugin.Start(doAll)
}
