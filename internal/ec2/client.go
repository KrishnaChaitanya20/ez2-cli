package ec2client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var ec2Client *ec2.Client

type Ec2Instances struct {
	InstanceID   string
	InstanceType string
	PublicIP     string
	PrivateIP    string
	Status       types.InstanceStateName
	InstanceName string
	Uptime       string
	// Cost         string
}

func InitClient(cfg aws.Config) {
	ec2Client = ec2.NewFromConfig(cfg)
}

func GetInstanceList(ctx context.Context, args []string) []Ec2Instances {
	// from cli extract args to filter tags
	input := ec2.DescribeInstancesInput{}
	input.Filters = processTags(args)

	// call the describe instances method
	resp, err := ec2Client.DescribeInstances(ctx, &input)
	if err != nil {
		log.Fatal(err)
	}

	return mapToEc2Instances(resp.Reservations)

}

func ChangeInstanceState(ctx context.Context, subcmd string, args []string) error {

	if len(args) == 0 {
		fmt.Println("Please provide instance id/name")
		fmt.Println("usage: ez2cli <start|stop|restart> id=<instance-id> [name=<instance-name> ...]")
		return nil
	}

	var ids []string
	var instanceIDRegex = regexp.MustCompile(`^i-[a-z0-9]{17}$`)

	for _, arg := range args {
		filterKey, filterVal, found := strings.Cut(arg, "=")

		if !found || filterKey == "" || filterVal == "" {
			fmt.Printf("Ignoring Invalid argument: %s\n", arg)
			continue
		}

		switch filterKey {
		case "id":
			if !instanceIDRegex.MatchString(filterVal) {
				fmt.Printf("Ignoring Invalid Instance ID: %s\n", filterVal)
				continue
			}
			ids = append(ids, filterVal)
		case "name":
			i, err := fetchInstanceFromInstanceName(ctx, filterVal)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			ids = append(ids, aws.ToString(i.InstanceId))
		}
	}

	if len(ids) == 0 {
		return errors.New("no valid instance id/name provided")
	}

	switch subcmd {
	case "start":
		output, err := ec2Client.StartInstances(ctx, &ec2.StartInstancesInput{InstanceIds: ids})
		if err != nil {
			return err
		}
		printStateTransition(output.StartingInstances)

	case "stop":
		output, err := ec2Client.StopInstances(ctx, &ec2.StopInstancesInput{InstanceIds: ids})
		if err != nil {
			return err
		}
		printStateTransition(output.StoppingInstances)

	case "restart":
		_, err := ec2Client.RebootInstances(ctx, &ec2.RebootInstancesInput{InstanceIds: ids})
		if err != nil {
			return err
		}
		fmt.Println("Rebooting Instances:", ids)

	}
	return nil
}

func ConnectToInstance(ctx context.Context, idOrName string, user string) (string, error) {
	var instanceIDRegex = regexp.MustCompile(`^i-[a-z0-9]{17}$`)
	isId := instanceIDRegex.MatchString(idOrName)

	var instance types.Instance
	var err error

	if !isId {
		instance, err = fetchInstanceFromInstanceName(ctx, idOrName)
	} else {
		instance, err = fetchInstanceFromInstanceId(ctx, idOrName)
	}

	if err != nil {
		log.Fatal(err)
	}

	if instance.State.Name != types.InstanceStateNameRunning {
		return "", errors.New("instance not running")
	}

	cmd := fmt.Sprintf(
		"ssh -i %s.pem %s@%s",
		*instance.KeyName,
		user,
		*instance.PublicIpAddress,
	)

	return cmd, nil
}

func mapToEc2Instances(reservations []types.Reservation) []Ec2Instances {
	var instances []Ec2Instances
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			var upTimeStr string
			if instance.State.Name == types.InstanceStateNameRunning {
				uptime := time.Since(*instance.LaunchTime)
				upTimeStr = fmt.Sprintf("%.0fh%.0fm", uptime.Hours(), uptime.Minutes())
			} else {
				upTimeStr = "N/A"
			}
			instances = append(instances, Ec2Instances{
				InstanceID:   aws.ToString(instance.InstanceId),
				InstanceType: string(instance.InstanceType),
				InstanceName: getTagValue(instance.Tags, "Name"),
				PrivateIP:    aws.ToString(instance.PrivateIpAddress),
				PublicIP:     aws.ToString(instance.PublicIpAddress),
				Status:       instance.State.Name,
				// Cost:         cost,
				Uptime: upTimeStr,
			})
		}
	}

	return instances
}

func fetchInstanceFromInstanceName(ctx context.Context, instanceName string) (types.Instance, error) {
	input := ec2.DescribeInstancesInput{}
	input.Filters = []types.Filter{
		{
			Name:   aws.String("tag:Name"),
			Values: []string{"*" + instanceName + "*"},
		},
	}
	resp, err := ec2Client.DescribeInstances(ctx, &input)
	if err != nil {
		return types.Instance{}, err
	}
	if len(resp.Reservations) == 0 {
		return types.Instance{}, errors.New("no Instances Found.")
	}
	if len(resp.Reservations) > 1 || len(resp.Reservations[0].Instances) > 1 {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{instanceName},
			},
		}
		resp, err = ec2Client.DescribeInstances(ctx, &input)
	}
	if err != nil {
		return types.Instance{}, err
	}

	return resp.Reservations[0].Instances[0], nil
}

func fetchInstanceFromInstanceId(ctx context.Context, instanceId string) (types.Instance, error) {
	input := ec2.DescribeInstancesInput{}
	input.InstanceIds = []string{instanceId}
	resp, err := ec2Client.DescribeInstances(ctx, &input)
	if err != nil {
		return types.Instance{}, err
	}
	if len(resp.Reservations) == 0 {
		return types.Instance{}, errors.New("no Instances Found.")
	}
	return resp.Reservations[0].Instances[0], nil
}

func printStateTransition(transitioningInstance []types.InstanceStateChange) {
	for _, instance := range transitioningInstance {
		fmt.Printf("Instance %s is now %s\n", *instance.InstanceId, instance.CurrentState.Name)
	}
}
