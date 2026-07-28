package main

import (
	"context"
	"ez2/utils"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var ec2Client *ec2.Client
var ctx context.Context

func main() {
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]
	var subcmd string
	if len(args) == 0 {
		fmt.Println("No subcommand provided")
		fmt.Println("Available subcommands: list, ls, stop, start")
		fmt.Println("Usage: ez2cli [subcommand] [args...]")
		return
	}
	subcmd = args[0]
	ec2Client = ec2.NewFromConfig(cfg)

	switch subcmd {
	case "ls", "list":
		instances := getInstanceList(args[1:])
		utils.PrintInstancesTable(instances)
	case "stop":
		stopInstance(args[1:])
	case "start":
		startInstance(args[1:])

	default:
		fmt.Println("Invalid usage ")
		fmt.Println("Available subcommands: list,ls,stop,start")
		fmt.Println("Usage: ez2cli [subcommand] [args...]")
		fmt.Println("ez2cli ls name=<tag:Name> | ez2cli ls id=<instance-id> | ez2cli ls state-<instance-state>")
		fmt.Println("ez2cli list name=<tag:Name> | ez2cli list id=<instance-id> | ez2cli list state-<instance-state>")
		fmt.Println("ez2cli start id|name <id|name> [<id2|name2> ...]")
		fmt.Println("ez2cli stop id|name <id|name> [<id2|name2> ...]")

	}

}

func getInstanceList(args []string) []utils.Ec2Instances {

	// from cli extract args to filter tags
	input := ec2.DescribeInstancesInput{}
	filterTags := utils.ProcessTags(args)
	input.Filters = filterTags

	// call the describe instances method
	resp, err := ec2Client.DescribeInstances(ctx, &input)
	if err != nil {
		log.Fatal(err)
	}
	var instances []utils.Ec2Instances

	// extract requried fields only
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, utils.Ec2Instances{
				InstanceID:   aws.ToString(instance.InstanceId),
				InstanceType: string(instance.InstanceType),
				InstanceName: utils.GetNameTag(instance.Tags),
				PrivateIP:    aws.ToString(instance.PrivateIpAddress),
				PublicIP:     aws.ToString(instance.PublicIpAddress),
				Status:       instance.State.Name,
			})
		}
	}
	return instances
}

func startInstance(args []string) {
	var input ec2.StartInstancesInput
	if len(args) == 0 {
		fmt.Println("No instance id provided")
		return
	}
	if args[0] == "id" {
		input = ec2.StartInstancesInput{
			InstanceIds: args[1:],
		}
	} else {
		fmt.Println("usage: ez2cli start id <id1> [<id2> ...]")
		return
	}
	output, err := ec2Client.StartInstances(ctx, &input)

	if err != nil {
		log.Fatal(err)
	}

	for _, instance := range output.StartingInstances {
		fmt.Printf("Instance %s is %s from %s \n", aws.ToString(instance.InstanceId), instance.CurrentState.Name, instance.PreviousState.Name)
	}
}

func stopInstance(args []string) {
	var input ec2.StopInstancesInput
	if len(args) == 0 {
		fmt.Println("No instance id provided")
		return
	}
	if args[0] == "id" {
		input = ec2.StopInstancesInput{
			InstanceIds: args[1:],
		}
	} else {
		fmt.Println("usage: ez2cli start id <id1> [<id2> ...]")
		return
	}
	output, err := ec2Client.StopInstances(ctx, &input)

	if err != nil {
		log.Fatal(err)
	}

	for _, instance := range output.StoppingInstances {
		fmt.Printf("Instance %s is %s from %s \n", aws.ToString(instance.InstanceId), instance.CurrentState.Name, instance.PreviousState.Name)
	}
}
