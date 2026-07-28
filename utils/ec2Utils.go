package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Ec2Instances struct {
	InstanceID   string
	InstanceType string
	PublicIP     string
	PrivateIP    string
	Status       types.InstanceStateName
	InstanceName string
}

func GetNameTag(tags []types.Tag) string {

	for _, tag := range tags {
		if *tag.Key == "Name" {
			return aws.ToString(tag.Value)
		}
	}

	return ""
}

func ProcessTags(args []string) []types.Filter {
	var filters []types.Filter

	for _, arg := range args {
		filterKey, filterVal, found := strings.Cut(arg, "=")

		if !found {
			fmt.Printf("Ignoring Invalid argument: %s\n", arg)
			continue
		}

		switch filterKey {
		case "id":
			filterKey = "instance-id"
		case "name":
			filterKey = "tag:Name"
		case "state":
			filterKey = "instance-state-name"
		default:
			filterKey = ""
		}

		filter, isPresent := findIn(filters, filterKey)

		if isPresent {
			filter.Values = append(filter.Values, filterVal)
		} else {
			filters = append(filters, types.Filter{
				Name:   &filterKey,
				Values: []string{filterVal},
			})
		}
	}

	return filters
}

func findIn(filters []types.Filter, filterKey string) (*types.Filter, bool) {
	for _, filter := range filters {
		if *filter.Name == filterKey {
			return &filter, true
		}
	}
	return &types.Filter{}, false
}
