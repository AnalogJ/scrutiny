package main

import (
	"fmt"
	"github.com/analogj/scrutiny/collector/pkg/collector"
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/errors"
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

	config, err := config.Create()
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	//we're going to load the config file manually, since we need to validate it.
	err = config.ReadConfig("/scrutiny/config/collector.yaml") // Find and read the config file
	if _, ok := err.(errors.ConfigFileMissingError); ok {      // Handle errors reading the config file
		//ignore "could not find config file"
	} else if err != nil {
		os.Exit(1)
	}

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
		Name:     "scrutiny-collector-metrics",
		Usage:    "smartctl data collector for scrutiny",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			collectorMetrics := "AnalogJ/scrutiny/metrics"

			var versionInfo string
			if len(goos) > 0 && len(goarch) > 0 {
				versionInfo = fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)
			} else {
				versionInfo = fmt.Sprintf("dev-%s", version.VERSION)
			}

			subtitle := collectorMetrics + utils.LeftPad2Len(versionInfo, " ", 65-len(collectorMetrics))

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
				Usage: "Run the scrutiny smartctl metrics collector",
				Action: func(c *cli.Context) error {
					if c.IsSet("config") {
						err = config.ReadConfig(c.String("config")) // Find and read the config file
						if err != nil {                             // Handle errors reading the config file
							//ignore "could not find config file"
							fmt.Printf("Could not find config file at specified path: %s", c.String("config"))
							return err
						}
					}
					//override config with flags if set
					if c.IsSet("host-id") {
						config.Set("host.id", c.String("host-id")) // set/override the host-id using CLI.
					}

					if c.Bool("debug") {
						config.Set("log.level", "DEBUG")
					}

					if c.IsSet("log-file") {
						config.Set("log.file", c.String("log-file"))
					}

					if c.IsSet("api-endpoint") {
						config.Set("api.endpoint", c.String("api-endpoint"))
					}

					collectorLogger := logrus.WithFields(logrus.Fields{
						"type": "metrics",
					})

					if level, err := logrus.ParseLevel(config.GetString("log.level")); err == nil {
						logrus.SetLevel(level)
					} else {
						logrus.SetLevel(logrus.InfoLevel)
					}

					if config.IsSet("log.file") && len(config.GetString("log.file")) > 0 {
						logFile, err := os.OpenFile(config.GetString("log.file"), os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							logrus.Errorf("Failed to open log file %s for output: %s", config.IsSet("log.file"), err)
							return err
						}
						defer logFile.Close()
						logrus.SetOutput(io.MultiWriter(os.Stderr, logFile))
					}

					metricCollector, err := collector.CreateMetricsCollector(
						config,
						collectorLogger,
						config.GetString("api.endpoint"),
					)

					if err != nil {
						return err
					}

					return metricCollector.Run()
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "config",
						Usage: "Specify the path to the devices file",
					},
					&cli.StringFlag{
						Name:    "api-endpoint",
						Usage:   "The api server endpoint",
						EnvVars: []string{"SCRUTINY_API_ENDPOINT"},
					},

					&cli.StringFlag{
						Name:    "log-file",
						Usage:   "Path to file for logging. Leave empty to use STDOUT",
						EnvVars: []string{"COLLECTOR_LOG_FILE"},
					},

					&cli.BoolFlag{
						Name:    "debug",
						Usage:   "Enable debug logging",
						EnvVars: []string{"COLLECTOR_DEBUG", "DEBUG"},
					},

					&cli.StringFlag{
						Name:    "host-id",
						Usage:   "Host identifier/label, used for grouping devices",
						Value:   "",
						EnvVars: []string{"COLLECTOR_HOST_ID"},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}

}
