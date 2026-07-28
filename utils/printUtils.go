package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintInstancesTable(instances []Ec2Instances) {

	if len(instances) == 0 {
		fmt.Println("No instances found.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Instance ID\tInstance Type\tState\tPublic IP\tPrivate IP\tName")
	fmt.Fprintln(w, "-----------\t-------------\t-----\t---------\t----------\t----")

	for _, instance := range instances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			instance.InstanceID,
			instance.InstanceType,
			instance.Status,
			instance.PublicIP,
			instance.PrivateIP,
			instance.InstanceName)
	}

	w.Flush()
}
