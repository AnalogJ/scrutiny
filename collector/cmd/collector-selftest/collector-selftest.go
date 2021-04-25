package main

import (
	"fmt"
	"github.com/analogj/scrutiny/collector/pkg/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/version"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"time"

	utils "github.com/analogj/go-util/utils"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var goos string
var goarch string

func main() {

	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}
CATEGORY:
   {{.Category}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	app := &cli.App{
		Name:     "scrutiny-collector-selftest",
		Usage:    "smartctl self-test data collector for scrutiny",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			collectorSelfTest := "AnalogJ/scrutiny/selftest"

			var versionInfo string
			if len(goos) > 0 && len(goarch) > 0 {
				versionInfo = fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)
			} else {
				versionInfo = fmt.Sprintf("dev-%s", version.VERSION)
			}

			subtitle := collectorSelfTest + utils.LeftPad2Len(versionInfo, " ", 65-len(collectorSelfTest))

			color.New(color.FgGreen).Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
				`
			 ___   ___  ____  __  __  ____  ____  _  _  _  _
			/ __) / __)(  _ \(  )(  )(_  _)(_  _)( \( )( \/ )
			\__ \( (__  )   / )(__)(   )(   _)(_  )  (  \  /
			(___/ \___)(_)\_)(______) (__) (____)(_)\_) (__)
			%s

			`), subtitle))

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the scrutiny self-test data collector",
				Action: func(c *cli.Context) error {

					collectorLogger := logrus.WithFields(logrus.Fields{
						"type": "selftest",
					})

					if c.Bool("debug") {
						logrus.SetLevel(logrus.DebugLevel)
					} else {
						logrus.SetLevel(logrus.InfoLevel)
					}

					if c.IsSet("log-file") {
						logFile, err := os.OpenFile(c.String("log-file"), os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							logrus.Errorf("Failed to open log file %s for output: %s", c.String("log-file"), err)
							return err
						}
						defer logFile.Close()
						logrus.SetOutput(io.MultiWriter(os.Stderr, logFile))
					}

					//TODO: pass in the collector, use configuration from collector-metrics
					stCollector, err := collector.CreateSelfTestCollector(
						collectorLogger,
						c.String("api-endpoint"),
					)

					if err != nil {
						return err
					}

					return stCollector.Run()
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "api-endpoint",
						Usage:   "The api server endpoint",
						Value:   "http://localhost:8080",
						EnvVars: []string{"SCRUTINY_API_ENDPOINT"},
					},

					&cli.StringFlag{
						Name:    "log-file",
						Usage:   "Path to file for logging. Leave empty to use STDOUT",
						Value:   "",
						EnvVars: []string{"COLLECTOR_LOG_FILE"},
					},

					&cli.BoolFlag{
						Name:    "debug",
						Usage:   "Enable debug logging",
						EnvVars: []string{"COLLECTOR_DEBUG", "DEBUG"},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}

}
