package cmds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
)

type ClusterConfig struct {
	project     string
	application string
}

var (
	clusterConfig     ClusterConfig
	autoscalingClient *autoscaling.AutoScaling
	clusterCommand    = &cobra.Command{
		Use:   "cluster",
		Short: "cluster management",
		Long:  ``,
	}
)

func init() {
	clusterCommand.PersistentFlags().StringVarP(&clusterConfig.application, "application", "a", "", "Application name")
	clusterCommand.PersistentFlags().StringVarP(&clusterConfig.project, "project", "j", "", "Project name")

	clusterCommand.AddCommand(listCommand)
	clusterCommand.AddCommand(describeCommand)
}

func getAutoscalingClient(parameters *CommandParameters) *autoscaling.AutoScaling {
	creds := credentials.NewCredentials(&credentials.SharedCredentialsProvider{
		Profile: parameters.profile,
	})

	if autoscalingClient == nil {
		autoscalingClient = autoscaling.New(session.New(), &aws.Config{
			Credentials: creds,
			Region:      aws.String(parameters.region),
		})
	}

	return autoscalingClient
}
