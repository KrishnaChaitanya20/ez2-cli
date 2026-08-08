package ec2client

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func processTags(args []string) []types.Filter {
	var filters []types.Filter

	for _, arg := range args {
		filterKey, filterVal, found := strings.Cut(arg, "=")

		if !found || filterKey == "" || filterVal == "" {
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

		filter, isPresent := findIn(&filters, filterKey)
		if isPresent {
			filter.Values = append(filter.Values[:], filterVal)
		} else {
			filters = append(filters, types.Filter{
				Name:   &filterKey,
				Values: []string{filterVal},
			})
		}
	}

	return filters
}

func findIn(filters *[]types.Filter, filterKey string) (*types.Filter, bool) {
	for i, filter := range *filters {
		if *filter.Name == filterKey {
			return &(*filters)[i], true
		}
	}
	return &types.Filter{}, false
}

func getTagValue(tags []types.Tag, tagName string) string {

	for _, tag := range tags {
		if *tag.Key == tagName {
			return aws.ToString(tag.Value)
		}
	}

	return ""
}

func extractFlags(args []string) []string {
	var flags []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			flags = append(flags, arg)
		}
	}
	return flags
}
