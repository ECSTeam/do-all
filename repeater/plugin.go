package repeater

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/xchapter7x/lo"
)

const name = "do-all"
const placeholder = "{}"

// Repeater the plugin struct that will be used for plugin executions
type Repeater struct {
	Version string
	Writer  io.Writer
}

// Run execute the plugin
func (c *Repeater) Run(cliConnection plugin.CliConnection, args []string) {
	args = args[1:]

	idx := -1
	for i, arg := range args {
		if arg == placeholder {
			idx = i
			break
		}
	}
	/*
		currentSpace, err := cliConnection.GetCurrentSpace()
		if err != nil {
			lo.G.Panic("PLUGIN ERROR: Could not determine current space: ", err)
			return
		} */

	apps, err := cliConnection.GetApps()
	if err != nil {
		lo.G.Panic("PLUGIN ERROR: get apps: ", err)
		return
	}

	for _, app := range apps {
		if idx >= 0 {
			args[idx] = app.Name
		}

		cliConnection.CliCommand(args...)
	}
}

// GetMetadata Return necessary metadata about the plugin
func (c *Repeater) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    name,
		Version: c.GetVersionType(),
		Commands: []plugin.Command{
			plugin.Command{
				Name:     name,
				HelpText: "Run the identified command on every app in a space. If the app name is a parameter in the command, use '{}'",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s scale {} -i 2", name),
				},
			},
		},
	}
}

// GetVersionType convert the semver string to a VersionType object
func (c *Repeater) GetVersionType() plugin.VersionType {
	versionArray := strings.Split(strings.TrimPrefix(c.Version, "v"), ".")
	major, _ := strconv.Atoi(versionArray[0])
	minor, _ := strconv.Atoi(versionArray[1])
	build, _ := strconv.Atoi(versionArray[2])
	return plugin.VersionType{
		Major: major,
		Minor: minor,
		Build: build,
	}
}
