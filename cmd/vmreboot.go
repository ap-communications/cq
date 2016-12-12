package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var vmrebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot running status VM",
	Long:  "Reboot running status VM",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup     //parallel processing counter group
		stats := map[string]int{} //instance id hit check map

		for _, argid := range args { //create hit judgment map for character string set as argument
			stats[argid] = 0 //init map (hit is 0)
		}

		if len(args) == 0 { //If there is no argument, abort
			fmt.Printf("missing args (Instance-ID)\n")
			return
		}

		regionsAWS := getAWSRegions()
		for _, region := range regionsAWS {
			wg.Add(1) //waiting group count up
			go rebootInstance(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {

	USAGE := `Usage:
  cq vm start [instance-id] [instance-id] ...
`

	vmCmd.AddCommand(vmrebootCmd)
	vmrebootCmd.SetUsageTemplate(USAGE)
}

func rebootInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	instanceParamEC2 := getEC2Param(region)

	for _, Reservations := range instanceParamEC2.Reservations {
		for _, Instances := range Reservations.Instances {

			for _, iid := range target {
				if *Instances.InstanceId == iid {
					if *Instances.State.Name == "running" { //if instance state is "running" (coz: if instance state is not "running" can't be reboot)
						ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})                            //generate API query instance
						_, err := ec2instance.RebootInstances(&ec2.RebootInstancesInput{InstanceIds: []*string{aws.String(iid)}}) //reboot instance
						if err != nil {
							fmt.Println(err)
							return
						}
						fmt.Printf("Success!   %s   has started reboot sequence\n", iid)
						stats[iid]++ //increment id hit counter
					} else {
						fmt.Printf("Can't reboot, Instance   %s   status is %s\n", iid, *Instances.State.Name) //if instance state is not "running" can't be reboot
						stats[iid]++                                                                           //increment id hit counter
					}
				}
			}

		}
	}

	return
}
