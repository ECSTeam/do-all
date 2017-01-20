package repeater_test

import (
	"bytes"
	"fmt"
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

			BeforeEach(func() {
				b = new(bytes.Buffer)
				plgn = &Repeater{
					Writer: b,
				}

				cli = new(pluginfakes.FakeCliConnection)

				fake_apps := make([]plugin_models.GetAppsModel, 0, 4)
				for _, app := range []string{"app1", "app2", "app3", "app4"} {
					fake_apps = append(fake_apps, plugin_models.GetAppsModel{
						Name: app,
					})
				}

				cli.GetAppsReturns(fake_apps, nil)

				cli.CliCommandStub = func(args ...string) ([]string, error) {
					s := fmt.Sprintf("do-all cf %v", strings.Join(args, " "))
					lines := make([]string, 0, 1)
					lines = append(lines, s)

					return lines, nil
				}
			})
			It("then it should list the commands it ran", func() {
				//	cli := new(pluginfakes.FakeCliConnection)

				plgn.Run(cli, []string{"do-all", "app", "{}"})

				out := b.String()

				for _, app := range []string{"app1", "app2", "app3", "app4"} {
					Ω(out).Should(ContainSubstring(fmt.Sprintf("do-all cf app %v", app)))
				}
			})
			// It("then it should print something", func() {
			// 	cli := new(pluginfakes.FakeCliConnection)
			// 	plgn.Run(cli, []string{"hi", "there"})
			// 	Ω(out).ShouldNot(BeEmpty())
			// })
		})
	})
})
