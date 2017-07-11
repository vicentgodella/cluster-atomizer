package cmds

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
	"github.com/vicentgodella/cluster-atomizer/pkg/repository"
	"testing"
)

func TestAutoscalingGroupsResultWriterWithNoGroups(t *testing.T) {
	groups := []*autoscaling.Group{}
	buf := new(bytes.Buffer)

	AutoscalingGroupsResultWriter(buf, groups)

	assert.Equal(t, "No autoscaling groups found", buf.String(), "It should be no autoscalings in result")
}

func TestAutoscalingGroupsResultWriterWithOneGroup(t *testing.T) {
	buf := new(bytes.Buffer)
	groups := []*autoscaling.Group{&autoscaling.Group{
		AutoScalingGroupName: aws.String("myproject-app-v000"),
		LoadBalancerNames: []*string{
			aws.String("myelb"),
		},
	}}

	AutoscalingGroupsResultWriter(buf, groups)

	expectedOutput := `+--------------------+-----------------+---------------------+
|        NAME        | INSTANCES COUNT | LOAD BALANCERS LIST |
+--------------------+-----------------+---------------------+
| myproject-app-v000 |               0 | myelb               |
+--------------------+-----------------+---------------------+
`
	assert.Equal(t, expectedOutput, buf.String())
}

func TestExecuteList(t *testing.T) {
	buf := new(bytes.Buffer)

	groups := []*autoscaling.Group{&autoscaling.Group{
		AutoScalingGroupName: aws.String("myproject-app-v000"),
		LoadBalancerNames: []*string{
			aws.String("myelb"),
		},
	}}
	repositoryMock := &inMemoryRepository{groups, "^myproject[\\-\\_]myapp\\-v([0-9]{3})$"}

	executor := ListClusterCommandExecutor{
		buf,
		repositoryMock,
		&ClusterConfig{"myproject", "myapp"},
	}

	executor.Execute()

	expectedOutput := `+--------------------+-----------------+---------------------+
|        NAME        | INSTANCES COUNT | LOAD BALANCERS LIST |
+--------------------+-----------------+---------------------+
| myproject-app-v000 |               0 | myelb               |
+--------------------+-----------------+---------------------+
`
	assert.Equal(t, expectedOutput, buf.String())
}

type inMemoryRepository struct {
	groups          []*autoscaling.Group
	expectedPattern string
}

func (m *inMemoryRepository) GetAutoscalingGroupsWithFilter(filter *repository.AutoscalingGroupsFilter) (error, []*autoscaling.Group) {
	if filter.Pattern != m.expectedPattern {
		return fmt.Errorf("Invalid pattern in filter. Actual: %s", filter.Pattern), nil
	}

	return nil, m.groups
}
