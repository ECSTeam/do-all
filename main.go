package main

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/do-all/repeater"
)

func main() {
	doAll := &repeater.Repeater{}

	plugin.Start(doAll)
}
