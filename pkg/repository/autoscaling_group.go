package repository

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"regexp"
)

type AutoscalingGroupRepository interface {
	GetAutoscalingGroupsWithFilter(filter *AutoscalingGroupsFilter) (error, []*autoscaling.Group)
}

// It is needed to be able to mock autoscaling.Autoscaling client.
// I have not found any interface for it in aws sdk.
type AutoscalingClient interface {
	DescribeAutoScalingGroups(*autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}

type RemoteAutoscalingGroupRepository struct {
	client AutoscalingClient
}

func (repo *RemoteAutoscalingGroupRepository) GetAutoscalingGroupsWithFilter(filter *AutoscalingGroupsFilter) (error, []*autoscaling.Group) {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}

	resp, err := repo.client.DescribeAutoScalingGroups(params)

	if err != nil {
		log.Error(err)
		return err, nil
	}

	return nil, FilterGroups(
		resp.AutoScalingGroups,
		filter,
	)
}

func NewAutoscalingGroupRepository(client *autoscaling.AutoScaling) AutoscalingGroupRepository {
	repository := &RemoteAutoscalingGroupRepository{
		client: client,
	}

	return repository
}

type AutoscalingGroupsFilter struct {
	Pattern string
}

func (f *AutoscalingGroupsFilter) filter(group *autoscaling.Group) bool {
	matches, err := regexp.Match(f.Pattern, []byte(*group.AutoScalingGroupName))

	if err != nil {
		return false
	}

	return matches
}

func FilterGroups(groups []*autoscaling.Group, filter *AutoscalingGroupsFilter) []*autoscaling.Group {
	var filteredGroups []*autoscaling.Group

	for _, v := range groups {
		if filter.filter(v) {
			filteredGroups = append(filteredGroups, v)
		}
	}

	return filteredGroups
}
