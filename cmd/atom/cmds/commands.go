package cmds

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

type CommandParameters struct {
	profile string
	region  string
}

var (
	commandConfig CommandParameters
	RootCommand   = &cobra.Command{
		Use:   "cluster-atom",
		Short: "Tool for deploying immutable clusters on cloud",
		Long:  ``,
	}
)

func init() {
	RootCommand.PersistentFlags().StringVarP(&commandConfig.profile,
		"profile",
		"p",
		"",
		"Profile name")
	RootCommand.PersistentFlags().StringVarP(&commandConfig.region,
		"region",
		"r",
		GetEnvOrDefault("AWS_DEFAULT_REGION", "eu-west-1"),
		"Region name")

	RootCommand.AddCommand(clusterCommand)
}

func Execute() {
	log.SetLevel(log.InfoLevel)

	if c, err := RootCommand.ExecuteC(); err != nil {
		if err != nil {
			c.Println(err.Error())
		}

		os.Exit(-1)
	}
}

func GetEnvOrDefault(var_name string, default_value string) string {
	value := os.Getenv(var_name)

	if value == "" {
		value = default_value
	}

	return value
}
