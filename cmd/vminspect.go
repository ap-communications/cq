package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var vminspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Print detailed instance info",
	Long:  "Print detailed instance info",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup                                  //parallel processing counter group
		var instanceParamEC2 = []ec2.DescribeInstancesOutput{} //API response slide

		stats := map[string]int{}    //instance id hit check map
		for _, argid := range args { //create hit judgment map for character string set as argument
			stats[argid] = 0 //init map (hit count is 0)
		}

		regionsAWS := getAWSRegions() //get region list (AWS)
		for _, region := range regionsAWS {
			wg.Add(1)                                      //waiting group count up
			go setEC2Param(region, &wg, &instanceParamEC2) //get instance information from API
			time.Sleep(1 * time.Millisecond)               //
		}
		wg.Wait()

		printInspect(instanceParamEC2, args, stats)
		printHitId(args, stats)

	},
}

func init() {

	USAGE := `Usage:
  cq vm inspect
    or
  cq vm inspect [instance-id] [instance-id] ...
`

	vmCmd.AddCommand(vminspectCmd)
	vminspectCmd.SetUsageTemplate(USAGE) //override Usage words
}

func printInspect(instanceParamEC2 []ec2.DescribeInstancesOutput, args []string, stats map[string]int) {

	if len(args) == 0 { //If there is no argument, print all instances information
		fmt.Println(instanceParamEC2)
	} else { //if there is argument instance id, print it
		for _, lawJson := range instanceParamEC2 {
			for _, Reservations := range lawJson.Reservations {
				for _, Instances := range Reservations.Instances {
					for _, iid := range args {
						if *Instances.InstanceId == iid {
							fmt.Println(*Instances)
							stats[iid]++ //increment hit id counter
						}
					}
				}
			}
		}
	}

	return

}
