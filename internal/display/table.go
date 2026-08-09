package display

import (
	"fmt"
	"os"
	"text/tabwriter"

	ec2client "ez2/internal/ec2"
)

func PrintInstancesTable(instances []ec2client.Ec2Instances) {

	if len(instances) == 0 {
		fmt.Println("No instances found.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Instance ID\tName\tInstance Type\tPublic IP\tPrivate IP\tState\tUptime")
	fmt.Fprintln(w, "-----------\t----\t-------------\t---------\t----------\t-----\t------")

	for _, instance := range instances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			instance.InstanceID,
			instance.InstanceName,
			instance.InstanceType,
			instance.PublicIP,
			instance.PrivateIP,
			instance.Status,
			instance.Uptime,
		)
	}

	w.Flush()
}
