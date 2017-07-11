package cmds

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
	"github.com/pkg/errors"
	"strings"
	"github.com/olekukonko/tablewriter"
	"bytes"
)

var describeCommand = &cobra.Command{
Use:   "describe",
Short: "describe cluster",
Long:  ``,
}
func init() {
}

func Map(vs []*string) string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] =  *v
	}

	return strings.Join(vsm, ",")
}

func dumpAutoscalingGroupInfo(group *autoscaling.Group, launchConfiguration *autoscaling.LaunchConfiguration) (string) {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)

	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string {
		"Name",
		*group.AutoScalingGroupName,
	})

	table.Append([]string {
		"ARN",
		*group.AutoScalingGroupARN,
	})

	var instancesList [] string
	for _, instance := range group.Instances {
		instancesList = append(instancesList, *instance.InstanceId)
	}
	table.Append([]string {
		"Instances",
		strings.Join(instancesList, ","),
	})

	table.Append([]string {
		"CreatedTime",
		(*group.CreatedTime).Format("2006-01-02 15:04:05 MST"),
	})
	table.Append([]string {
		"MinSize",
		fmt.Sprint(*group.MinSize),

	})
	table.Append([]string {
		"DesiredCapacity",
		fmt.Sprint(*group.DesiredCapacity),
	})
	table.Append([]string {
		"MaxSize",
		fmt.Sprint(*group.MaxSize),
	})
	table.Append([]string {
		"ImageId",
		*launchConfiguration.ImageId,
	})
	table.Append([]string {
		"SecurityGroups",
		Map(launchConfiguration.SecurityGroups),
	})

	table.Render()

	return buf.String()
}

func getAutoscalingGroupInfo(client *autoscaling.AutoScaling, name string) (error, *autoscaling.Group) {
	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&name},
	}

	resp, err := client.DescribeAutoScalingGroups(params)

	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}

	if len(resp.AutoScalingGroups) != 1 {
		message := fmt.Sprintf("Autoscaling %s not found", name)
		err = errors.New(message)
		fmt.Println(err.Error())
		return err, nil
	}

	return nil, resp.AutoScalingGroups[0]
}

func getLaunchConfigurationInfo(client *autoscaling.AutoScaling, name string) (error, *autoscaling.LaunchConfiguration) {
	params := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{&name},
	}

	resp, err := client.DescribeLaunchConfigurations(params)

	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}

	if len(resp.LaunchConfigurations) != 1 {
		message := fmt.Sprintf("LaunchConfiguration %s not found", name)
		err = errors.New(message)
		fmt.Println(err.Error())
		return err, nil
	}

	return nil, resp.LaunchConfigurations[0]
}


func NewDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe <autoscaling-name>",
		Short: "Describe cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Autoscaling name is missing")
			}

			autoscalingName := strings.Join(args, " ")

			client := getAutoscalingClient(&commandConfig)
			err, group := getAutoscalingGroupInfo(client, autoscalingName)

			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			err, launchConfiguration := getLaunchConfigurationInfo(client, *group.LaunchConfigurationName)

			fmt.Println(launchConfiguration)
			fmt.Println(group)

			fmt.Println(dumpAutoscalingGroupInfo(group, launchConfiguration))

			return nil
		},
	}

	return cmd
}
