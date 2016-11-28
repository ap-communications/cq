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

var vmstartCmd = &cobra.Command{
	Use:   "start",
	Short: "Startup VM",
	Long:  "Startup VM",
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
			go startInstance(args, region, &wg, stats)
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

	vmCmd.AddCommand(vmstartCmd)
	vmstartCmd.SetUsageTemplate(USAGE) //override Usage words
}

func startInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	instanceParamEC2 := getEC2Param(region)

	for _, Reservations := range instanceParamEC2.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, iid := range target {
				if *Instances.InstanceId == iid {
					ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})                             //generate API query instance
					resp, err := ec2instance.StartInstances(&ec2.StartInstancesInput{InstanceIds: []*string{aws.String(iid)}}) //start instance
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("Success!  %s   %s  ===>  %s\n", iid, *resp.StartingInstances[0].PreviousState.Name, *resp.StartingInstances[0].CurrentState.Name)
					stats[iid]++ //increment id hit counter
				}
			}
		}
	}

	wg.Done()
	return

}
