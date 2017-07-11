package repository

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAutoscalingGroupsWithFilter(t *testing.T) {
	groups := []*autoscaling.Group{&autoscaling.Group{
		AutoScalingGroupName: aws.String("myproject-v000"),
		LoadBalancerNames: []*string{
			aws.String("myelb"),
		},
	}}
	autoscalingGroupClient := &inMemoryASGClient{
		groups,
		false,
	}

	repository := RemoteAutoscalingGroupRepository{autoscalingGroupClient}
	filter := &AutoscalingGroupsFilter{NewProjectFilterPatternBuilder("myproject").Build()}
	err, groups := repository.GetAutoscalingGroupsWithFilter(filter)

	assert.NotNil(t, groups)
	assert.Nil(t, err)
}
func TestGetAutoscalingGroupsWithFilterReturnsErrorIfSomethingBadHappens(t *testing.T) {
	autoscalingGroupClient := &inMemoryASGClient{
		nil,
		true,
	}

	repository := RemoteAutoscalingGroupRepository{autoscalingGroupClient}
	filter := &AutoscalingGroupsFilter{NewProjectFilterPatternBuilder("myproject").Build()}
	err, groups := repository.GetAutoscalingGroupsWithFilter(filter)

	assert.Nil(t, groups)
	assert.NotNil(t, err)
}

type inMemoryASGClient struct {
	groups            []*autoscaling.Group
	shouldReturnError bool
}

func (m *inMemoryASGClient) DescribeAutoScalingGroups(*autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	if m.shouldReturnError {
		return nil, errors.New("Fake error")
	}
	output := &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: m.groups,
	}

	return output, nil
}
