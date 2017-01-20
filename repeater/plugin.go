package repeater

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/plugin/models"

	"github.com/cloudfoundry/cli/plugin"

	"github.com/krujos/cfcurl"

	"github.com/xchapter7x/lo"
)

const name = "do-all"
const placeholder = "{}"

var version = "0.0.1"

// Repeater the plugin struct that will be used for plugin executions
type Repeater struct {
	Writer io.Writer
}

type orgInfo struct {
	name   string
	spaces []string
}

// Run execute the plugin
func (c *Repeater) Run(cli plugin.CliConnection, args []string) {
	var firstArg = args[1]

	var orgInfos []orgInfo
	var err error

	var currentOrg plugin_models.Organization
	var currentSpace plugin_models.Space

	defer func() {
		cli.CliCommand("target", "-o", currentOrg.Name, "-s", currentSpace.Name)
	}()

	if strings.HasPrefix(firstArg, "--") {
		if len(args) < 3 {
			fmt.Printf("You have to tell do-all to do something!")
			lo.G.Panic("You have to tell do-all to do something!")
		}

		args = args[2:]

		orgInfos, err = c.getAllOrgsAndSpaces(cli, (firstArg == "--global"))
	} else {
		if len(args) < 2 {
			fmt.Printf("You have to tell do-all to do something!")
			lo.G.Panic("You have to tell do-all to do something!")
		}

		orgInfos, err = c.getCurrentOrgAndSpace(cli)

		args = args[1:]
	}

	if err != nil {
		lo.G.Panic("PLUGIN ERROR: get apps: ", err)
		return
	}

	// capture current target
	currentOrg, err = cli.GetCurrentOrg()
	if err != nil {
		lo.G.Panic("PLUGIN ERROR: get apps: ", err)
		return
	}

	currentSpace, err = cli.GetCurrentSpace()
	if err != nil {
		lo.G.Panic("PLUGIN ERROR: get apps: ", err)
		return
	}

	for _, orgInfo := range orgInfos {
		for _, space := range orgInfo.spaces {
			c.runCommands(cli, orgInfo.name, space, args)
		}
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
					Usage: fmt.Sprintf("cf %s [--org|--global] scale {} -i 2", name),
					Options: map[string]string{
						"org":    "Run the identified command on every app in the org, instead of the space",
						"global": "Run the idenitified command on every app globally, instead of the space",
					},
				},
			},
		},
	}
}

// GetVersionType convert the semver string to a VersionType object
func (c *Repeater) GetVersionType() plugin.VersionType {
	versionArray := strings.Split(strings.TrimPrefix(version, "v"), ".")
	major, _ := strconv.Atoi(versionArray[0])
	minor, _ := strconv.Atoi(versionArray[1])
	build, _ := strconv.Atoi(versionArray[2])
	return plugin.VersionType{
		Major: major,
		Minor: minor,
		Build: build,
	}
}

func (c *Repeater) runCommands(cli plugin.CliConnection, orgName string, spaceName string, args []string) {
	idx := -1
	for i, arg := range args {
		if arg == placeholder {
			idx = i
			break
		}
	}

	cli.CliCommand("target", "-o", orgName, "-s", spaceName)
	apps, _ := cli.GetApps()

	for _, app := range apps {
		if idx >= 0 {
			args[idx] = app.Name
		}

		var cmdOutput []string
		var err error
		if cmdOutput, err = cli.CliCommand(args...); err != nil {
			lo.G.Panic(err)
		}

		if c.Writer != nil {
			for _, line := range cmdOutput {
				fmt.Fprint(c.Writer, line)
			}
		}
	}
}

func (c *Repeater) getCurrentOrgAndSpace(cli plugin.CliConnection) ([]orgInfo, error) {
	o, err := cli.GetCurrentOrg()
	if err != nil {
		return nil, err
	}

	s, err := cli.GetCurrentSpace()
	if err != nil {
		return nil, err
	}

	return []orgInfo{
		orgInfo{
			name:   o.Name,
			spaces: []string{s.Name},
		},
	}, nil
}

func (c *Repeater) getAllOrgsAndSpaces(cli plugin.CliConnection, global bool) ([]orgInfo, error) {
	var orgs []models.OrganizationFields
	var orgInfos []orgInfo

	if global {
		orgModels, err := cli.GetOrgs()
		if err != nil {
			return nil, err
		}

		orgs = make([]models.OrganizationFields, len(orgModels))
		for _, model := range orgModels {
			orgs = append(orgs, models.OrganizationFields{
				Name: model.Name,
				GUID: model.Guid,
			})
		}
	} else {
		currentOrg, err := cli.GetCurrentOrg()
		if err != nil {
			return nil, err
		}

		orgs = make([]models.OrganizationFields, 1)
		orgs = append(orgs, models.OrganizationFields{
			Name: currentOrg.Name,
			GUID: currentOrg.Guid,
		})
	}

	orgInfos = make([]orgInfo, len(orgs))
	for _, org := range orgs {
		var nextURL interface{}
		nextURL = "/v2/spaces?q=organization_guid+IN+" + org.GUID

		spaceNames := make([]string, 5)
		for nextURL != nil {
			json, err := cfcurl.Curl(cli, nextURL.(string))

			if err != nil {
				return nil, err
			}

			resources := toJSONArray(json["resources"])
			for _, spaceIntf := range resources {
				space := toJSONObject(spaceIntf)
				entity := toJSONObject(space["entity"])
				spaceNames = append(spaceNames, entity["name"].(string))
			}

			nextURL = json["next_url"]
		}

		orgInfos = append(orgInfos, orgInfo{
			name:   org.Name,
			spaces: spaceNames,
		})
	}

	return orgInfos, nil
}

func toJSONArray(obj interface{}) []interface{} {
	return obj.([]interface{})
}

func toJSONObject(obj interface{}) map[string]interface{} {
	return obj.(map[string]interface{})
}
