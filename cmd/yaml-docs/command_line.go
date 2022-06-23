package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andbron/yaml-docs/pkg/document"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string = "0.5.0"

func possibleLogLevels() []string {
	levels := make([]string, 0)

	for _, l := range log.AllLevels {
		levels = append(levels, l.String())
	}

	return levels
}

func initializeCli() {
	logLevelName := viper.GetString("log-level")
	logLevel, err := log.ParseLevel(logLevelName)
	if err != nil {
		log.Errorf("Failed to parse provided log level %s: %s", logLevelName, err)
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(logLevel)
}

func newYAMLDocsCommand(run func(cmd *cobra.Command, args []string)) (*cobra.Command, error) {
	command := &cobra.Command{
		Use:     "yaml-docs",
		Short:   "yaml-docs automatically generates markdown documentation from yaml values files",
		Version: version,
		Run:     run,
	}

	logLevelUsage := fmt.Sprintf("Level of logs that should printed, one of (%s)", strings.Join(possibleLogLevels(), ", "))
	command.PersistentFlags().BoolP("dry-run", "d", false, "don't actually render any markdown files just print to stdout passed")
	command.PersistentFlags().StringP("log-level", "l", "info", logLevelUsage)
	command.PersistentFlags().StringP("output-file", "o", "README.md", "markdown file path relative to input template to which rendered documentation will be written")
	command.PersistentFlags().StringP("sort-values-order", "s", document.AlphaNumSortOrder, fmt.Sprintf("order in which to sort the values table (\"%s\" or \"%s\")", document.AlphaNumSortOrder, document.FileSortOrder))
	command.PersistentFlags().StringSliceP("template-files", "t", []string{"README.md.gotmpl"}, "gotemplate file paths relative to each chart directory from which documentation will be generated")
	command.PersistentFlags().StringSliceP("values-file", "f", []string{}, "yaml values file to be parsed into values table. Can be specified multiple times")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("YAML_DOCS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.BindPFlags(command.PersistentFlags())

	return command, err
}
