package repeater_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin/models"

	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	. "github.com/ecsteam/do-all/repeater"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repeater", func() {
	Describe("Given a Repeater object", func() {
		Context("when calling Run w/ the wrong args", func() {
			var plgn *Repeater
			BeforeEach(func() {
				plgn = new(Repeater)
			})
			It("then it should panic and exit", func() {
				Ω(func() {
					cli := new(pluginfakes.FakeCliConnection)
					plgn.Run(cli, []string{})
				}).Should(Panic())
			})
		})

		Context("when calling Run w/ the proper args", func() {
			var plgn *Repeater
			var b *bytes.Buffer
			var cli *pluginfakes.FakeCliConnection

			var orgs []plugin_models.GetOrgs_Model

			var currentOrg plugin_models.Organization = plugin_models.Organization{
				plugin_models.OrganizationFields{
					Name: "ecsteam",
					Guid: "88195490-843f-4342-8e24-bf9a7bc690f4",
				},
			}

			var currentSpace plugin_models.Space = plugin_models.Space{
				plugin_models.SpaceFields{
					Name: "development",
				},
			}

			BeforeEach(func() {
				b = new(bytes.Buffer)
				plgn = &Repeater{
					Writer: b,
				}

				cli = new(pluginfakes.FakeCliConnection)

				cli.GetOrgsStub = func() ([]plugin_models.GetOrgs_Model, error) {
					file, _ := ioutil.ReadFile("fixtures/GetOrgs.json")

					if orgs == nil {
						json.Unmarshal(file, &orgs)
					}

					return orgs, nil
				}

				cli.GetCurrentOrgStub = func() (plugin_models.Organization, error) {
					return currentOrg, nil
				}

				cli.GetCurrentSpaceStub = func() (plugin_models.Space, error) {
					return currentSpace, nil
				}

				cli.GetAppsStub = func() ([]plugin_models.GetAppsModel, error) {
					file := strings.Join([]string{"fixtures", "orgs", currentOrg.Name, currentSpace.Name, "apps.json"}, "/")

					bytes, _ := ioutil.ReadFile(file)

					var apps []plugin_models.GetAppsModel
					json.Unmarshal(bytes, &apps)

					return apps, nil
				}

				setTarget := func(cli *pluginfakes.FakeCliConnection, args []string) {
					orgs, _ := cli.GetOrgs()

					for _, org := range orgs {
						if org.Name == args[2] {
							currentOrg = plugin_models.Organization{
								plugin_models.OrganizationFields{
									Name: org.Name,
									Guid: org.Guid,
								},
							}
							break
						}
					}

					currentSpace = plugin_models.Space{
						plugin_models.SpaceFields{
							Name: args[4],
						},
					}
				}

				cli.CliCommandWithoutTerminalOutputStub = func(args ...string) ([]string, error) {
					var output []string

					if args[0] == "target" {
						setTarget(cli, args)
						return []string{}, nil
					}

					r := strings.NewReplacer("?", "_", "=", "_", "+", "_")
					name := r.Replace(args[1])

					fixtureName := "fixtures" + name + ".json"

					file, err := os.Open(fixtureName)
					defer file.Close()
					if err != nil {
						Fail("Could not open " + fixtureName)
					}

					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						output = append(output, scanner.Text())
					}

					return output, scanner.Err()
				}

				cli.CliCommandStub = func(args ...string) ([]string, error) {
					if args[0] == "target" {
						setTarget(cli, args)
						return []string{}, nil
					}

					s := fmt.Sprintf("cf do-all %v\n", strings.Join(args, " "))
					lines := make([]string, 0, 1)
					lines = append(lines, s)

					return lines, nil
				}
			})
			It("then it should list the commands it ran in space only mode", func() {
				plgn.Run(cli, []string{"do-all", "app", "{}"})

				out := b.String()

				for _, app := range []string{"spring-music", "willitconnect"} {
					Ω(out).Should(ContainSubstring(fmt.Sprintf("cf do-all app %v\n", app)))
				}
			})

			It("then it should list the commands it ran in org mode", func() {
				plgn.Run(cli, []string{"do-all", "--org", "app", "{}"})

				out := b.String()

				for _, app := range []string{"spring-music",
					"willitconnect",
					"plugin-test-demo",
				} {
					Ω(out).Should(ContainSubstring(fmt.Sprintf("cf do-all app %v\n", app)))
				}
			})

			It("then it should list the commands it ran in global mode", func() {
				//	cli := new(pluginfakes.FakeCliConnection)

				plgn.Run(cli, []string{"do-all", "--global", "app", "{}"})

				out := b.String()

				for _, app := range []string{
					"Downloads/spring-music.war",
					"sp-music",
					"spring-joker/",
					"spring-music/",
					"spring-music",
					"willitconnect",
					"plugin-test-demo",
					"hello-goodbye",
					"helloworld-no-mvc",
					"helloworld-with-mvc",
					"mvc-demo02",
					"sp-music01",
					"clouddriver",
					"deck",
					"echo",
					"spinnaker_app",
					"autoscale",
					"azure-service-broker-1.2.0",
					"notifications-ui",
					"spring-cloud-broker",
					"spring-cloud-broker-worker",
					"pivotal-account",
					"pivotal-account-cold",
					"app-usage-scheduler",
					"app-usage-scheduler-venerable",
					"app-usage-server",
					"app-usage-server-venerable",
					"app-usage-worker",
					"app-usage-worker-venerable",
					"apps-manager-js",
					"apps-manager-js-venerable",
					"p-invitations",
					"p-invitations-venerable",
					"php-demo",
					"spring-music",
				} {
					Ω(out).Should(ContainSubstring(fmt.Sprintf("cf do-all app %v\n", app)))
				}
			})
		})
	})
})
