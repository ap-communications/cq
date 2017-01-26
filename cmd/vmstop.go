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

var vmstopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Shutdown VM",
	Long:  "Shutdown VM",
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

		if listFlag.Force { //if there is enabled force option, dont confirmation
			//jump to stop sequence
		} else {
			input := ""                                                              //keyboard input value
			fmt.Printf("Instance   %s   will be stop, are you sure?  Y/N\n", args) //destroy warning (pre)
			fmt.Scanln(&input)                                                       //stdin
			if (input == "Y") || (input == "y") {                                    //input Y or y
				//jump to stop sequence
			} else {
				fmt.Printf("Cancelled\n")
				return
			}
		}

		regionsAWS := getAWSRegions()
		for _, region := range regionsAWS {
			wg.Add(1) //waiting group count up
			go stopInstance(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond)
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {
	vmCmd.AddCommand(vmstopCmd)
	vmstopCmd.Flags().BoolVarP(&listFlag.Force, "force", "f", false, "Stop without confirmation") //define -f --force flag
}

func stopInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	instanceParamEC2 := getEC2Param(region)

	for _, Reservations := range instanceParamEC2.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, iid := range target {
				if *Instances.InstanceId == iid {
					ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})                           //generate API query instance
					resp, err := ec2instance.StopInstances(&ec2.StopInstancesInput{InstanceIds: []*string{aws.String(iid)}}) //stop instance
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("Success!  %s   %s  ===>  %s\n", iid, *resp.StoppingInstances[0].PreviousState.Name, *resp.StoppingInstances[0].CurrentState.Name)
					stats[iid]++ //increment id hit counter
				}
			}
		}
	}

	return

}
