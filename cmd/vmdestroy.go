package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"strings"
	"sync"
	"time"
)

var vmdestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "[DANGER]  Destroy VM  (CAN NOT RESTORE)",
	Long:  "[DANGER]  Destroy VM  (CAN NOT RESTORE)",
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

		ids := ""                      //keyboard input value
		for _, inputid := range args { //translate comma spreaded (for warning print)
			ids += inputid + ", "
		}
		ids = strings.TrimRight(ids, ", ") //delete final comma

		if listFlag.Force { //if there is enabled force option, dont confirmation
			regionsAWS := getAWSRegions()
			for _, region := range regionsAWS {
				wg.Add(1)                                    //waiting group count up
				go destroyInstance(args, region, &wg, stats) //destroy instance
				time.Sleep(1 * time.Millisecond)             //
			}
			wg.Wait()
		} else {
			input := ""                                                              //keyboard input value
			fmt.Printf("Instance   %s   will be DESTROY, are you sure?  Y/N\n", ids) //destroy warning (pre)
			fmt.Scanln(&input)                                                       //stdin
			if (input == "Y") || (input == "y") {                                    //input Y or y (1)
				fmt.Printf("This is final warning. DESTROY instance   %s   ARE YOU SURE? (Check EBS data)  Y/N\n", ids) //final destroy warning
				fmt.Scanln(&input)                                                                                      //stdin
				if (input == "Y") || (input == "y") {                                                                   //input Y or y (2 (final))
					regionsAWS := getAWSRegions()
					for _, region := range regionsAWS {
						wg.Add(1)                                    //waiting group count up
						go destroyInstance(args, region, &wg, stats) //destroy instance
						time.Sleep(1 * time.Millisecond)             //
					}
					wg.Wait()
				} else { //not Y or y, exit
					fmt.Printf("Cancelled\n")
					return
				}
			} else { //not Y or y, exit
				fmt.Printf("Cancelled\n")
				return
			}
		}

		printHitId(args, stats)

	},
}

func init() {
	vmCmd.AddCommand(vmdestroyCmd)
	vmdestroyCmd.Flags().BoolVarP(&listFlag.Force, "force", "f", false, "Destroy without confirmation") //define -f --force flag
}

func destroyInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	instanceParamEC2 := getEC2Param(region)

	for _, Reservations := range instanceParamEC2.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, iid := range target {
				if *Instances.InstanceId == iid {
					ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})                                  //generate API query instance
					_, err := ec2instance.TerminateInstances(&ec2.TerminateInstancesInput{InstanceIds: []*string{aws.String(iid)}}) //destroy instance
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("Success!   %s   destroyed\n", iid)
					stats[iid]++ //increment hit id counter
				}
			}
		}
	}

	wg.Done()
	return

}
