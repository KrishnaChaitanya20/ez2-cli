package ec2client

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type FormatFilter types.Filter

func (f FormatFilter) String() string {
	return fmt.Sprintf("%s=%v", *f.Name, f.Values)
}

func mapToMYFilter(filters []types.Filter) []FormatFilter {
	res := []FormatFilter{}
	for _, f := range filters {
		res = append(res, FormatFilter(f))
	}
	return res
}

func TestProcessTags(t *testing.T) {

	happyPathTestCases := [][]string{
		{"id=i-0fa3f8090975c39f8"},
		{"name=MC 26.2"},
		{"name=M*"},
		{"id=i-0fa3f8090975c39f8", "name=MC 26.2"},
		{"id=i-0fa3f8090975c39f8", "id=i-1fa3f8090975c39f8"},
		{"name=MC 26.2", "name=Prod"},
		{"name=MC 26.2", "name=Prod", "name=Prod*"},
		{"name="},
		{"status"},
		{"=MC"},
	}

	t.Run("Testing Happy paths", func(t *testing.T) {
		for _, test := range happyPathTestCases {
			output := processTags(test)
			fmt.Println("Test:", test, "\tOutput", mapToMYFilter(output))
		}
	})

}
