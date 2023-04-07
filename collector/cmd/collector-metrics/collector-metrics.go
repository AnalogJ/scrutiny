package main

import (
	"encoding/json"
	"fmt"
	"github.com/analogj/scrutiny/collector/pkg/collector"
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/webapp/backend/pkg/version"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"strings"
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

	configFilePath := "/opt/scrutiny/config/collector.yaml"
	configFilePathAlternative := "/opt/scrutiny/config/collector.yml"
	if !utils.FileExists(configFilePath) && utils.FileExists(configFilePathAlternative) {
		configFilePath = configFilePathAlternative
	}

	//we're going to load the config file manually, since we need to validate it.
	err = config.ReadConfig(configFilePath) // Find and read the config file
	if _, ok := err.(errors.ConfigFileMissingError); ok {          // Handle errors reading the config file
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
						//if the user is providing an api-endpoint with a basepath (eg. http://localhost:8080/scrutiny),
						//we need to ensure the basepath has a trailing slash, otherwise the url.Parse() path concatenation doesnt work.
						apiEndpoint := strings.TrimSuffix(c.String("api-endpoint"), "/") + "/"
						config.Set("api.endpoint", apiEndpoint)
					}

					collectorLogger, logFile, err := CreateLogger(config)
					if logFile != nil {
						defer logFile.Close()
					}
					if err != nil {
						return err
					}

					settingsData, err := json.MarshalIndent(config.AllSettings(), "", "\t")
					collectorLogger.Debug(string(settingsData), err)
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
						EnvVars: []string{"COLLECTOR_API_ENDPOINT", "SCRUTINY_API_ENDPOINT"},
						//SCRUTINY_API_ENDPOINT is deprecated, but kept for backwards compatibility
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

func CreateLogger(appConfig config.Interface) (*logrus.Entry, *os.File, error) {
	logger := logrus.WithFields(logrus.Fields{
		"type": "metrics",
	})

	if level, err := logrus.ParseLevel(appConfig.GetString("log.level")); err == nil {
		logger.Logger.SetLevel(level)
	} else {
		logger.Logger.SetLevel(logrus.InfoLevel)
	}

	var logFile *os.File
	var err error
	if appConfig.IsSet("log.file") && len(appConfig.GetString("log.file")) > 0 {
		logFile, err = os.OpenFile(appConfig.GetString("log.file"), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Logger.Errorf("Failed to open log file %s for output: %s", appConfig.GetString("log.file"), err)
			return nil, logFile, err
		}
		logger.Logger.SetOutput(io.MultiWriter(os.Stderr, logFile))
	}
	return logger, logFile, nil
}
