package cmds

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"

	log "github.com/Sirupsen/logrus"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vicentgodella/cluster-atomizer/pkg/repository"
	"io"
	"os"
	"strconv"
	"strings"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Long:  ``,
}

func init() {
	listCommand.RunE = ListCluster
}

func AutoscalingGroupsResultWriter(writer io.Writer, groups []*autoscaling.Group) {
	if len(groups) == 0 {
		writer.Write([]byte("No autoscaling groups found"))
		return
	}

	table := tablewriter.NewWriter(writer)

	for _, group := range groups {
		var loadBalancersList []string
		for _, name := range group.LoadBalancerNames {
			loadBalancersList = append(loadBalancersList, *name)
		}
		item := []string{
			*group.AutoScalingGroupName,
			strconv.Itoa(len(group.Instances)),
			strings.Join(loadBalancersList, ","),
		}
		table.Append(item)
	}

	table.SetHeader([]string{"Name", "Instances count", "Load Balancers List"})
	table.Render()
}

func ListCluster(cmd *cobra.Command, args []string) error {
	client := getAutoscalingClient(&commandConfig)
	asgRepository := repository.NewAutoscalingGroupRepository(client)

	writer := io.Writer(os.Stdout)
	executor := ListClusterCommandExecutor{
		writer,
		asgRepository,
		&clusterConfig,
	}
	return executor.Execute()
}

type ListClusterCommandExecutor struct {
	writer        io.Writer
	asgRepository repository.AutoscalingGroupRepository
	clusterConfig *ClusterConfig
}

func (ex *ListClusterCommandExecutor) Execute() error {
	filterPattern := ex.buildFilterPattern()
	filter := &repository.AutoscalingGroupsFilter{filterPattern}

	err, groups := ex.asgRepository.GetAutoscalingGroupsWithFilter(filter)

	if err != nil {
		log.Error(err)
		return err
	}

	AutoscalingGroupsResultWriter(ex.writer, groups)
	return nil
}
func (ex *ListClusterCommandExecutor) buildFilterPattern() string {
	var builder repository.FilterPatternBuilder
	builder = repository.NewProjectFilterPatternBuilder(ex.clusterConfig.project)

	if ex.clusterConfig.application != "" {
		builder = repository.NewProjectAndApplicationFilterPatternBuilder(
			ex.clusterConfig.project,
			ex.clusterConfig.application,
		)
	}
	return builder.Build()
}
