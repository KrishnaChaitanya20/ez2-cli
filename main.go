package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"ez2/internal/display"
	ec2client "ez2/internal/ec2"
	pricingclient "ez2/internal/pricing"
)

func main() {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal(err)
	}
	ec2client.InitClient(cfg)
	pricingclient.InitClient(cfg)

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No subcommand provided")
		fmt.Println("Available subcommands: list, ls, stop, start")
		fmt.Println("Usage: ez2cli <subcommand> [args...]")
		return
	}

	var subcmd string
	subcmd = args[0]

	switch subcmd {
	case "ls", "list":
		instances := ec2client.GetInstanceList(ctx, args[1:])
		display.PrintInstancesTable(instances)
	case "stop", "start", "restart":
		ec2client.ChangeInstanceState(ctx, subcmd, args[1:])
	case "connect":
		if len(args[1:]) != 2 {
			fmt.Println("Invalid Arguments")
			fmt.Println("usage: ez2 connect <instanceid/Name> <user>")
			return
		}
		idOrName := args[1]
		user := args[2]

		connectionString, err := ec2client.ConnectToInstance(ctx, idOrName, user)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(connectionString)

	default:
		fmt.Println("Invalid usage ")
		fmt.Println("Available subcommands: list,ls,stop,start")

	}

}
