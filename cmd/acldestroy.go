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

var acldestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy ACL (CAN NOT RESTORE)",
	Long: `destroy ACL (CAN NOT RESTORE)

Example:
  cq acl destroy sg-fd8cc1ee
`,
	Run: func(cmd *cobra.Command, args []string) {

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
			startParallelsDestroySecurityGroup(args, stats)
		} else {
			input := ""                                                                                    //keyboard input value
			fmt.Printf("SecurityGroup   %s   will be DESTROY, are you sure? (CAN NOT RESTORE) Y/N\n", ids) //destroy warning
			fmt.Scanln(&input)                                                                             //stdin
			if (input == "Y") || (input == "y") {                                                          //input Y or y
				startParallelsDestroySecurityGroup(args, stats)
			} else { //not Y or y, exit
				fmt.Printf("Cancelled\n")
				return
			}
		}

		printHitId(args, stats)

	},
}

func init() {
	aclCmd.AddCommand(acldestroyCmd)
	acldestroyCmd.Flags().BoolVarP(&listFlag.Force, "force", "f", false, "Destroy without confirmation") //define -f --force flag
}

func startParallelsDestroySecurityGroup(args []string, stats map[string]int) {

	var wg sync.WaitGroup     //parallel processing counter group

	regionsAWS := getAWSRegions()
	for _, region := range regionsAWS {
		wg.Add(1)                                    //waiting group count up
		go destroySecurityGroup(region, &wg, stats, args) //destroy instance
		time.Sleep(1 * time.Millisecond)             //
	}
	wg.Wait()

}

func destroySecurityGroup(region string, wg *sync.WaitGroup, stats map[string]int, target []string) {

	sgParamEC2 := getSecurityGroupParam(region) //get security group parameter

	for _, SecurityGroups := range sgParamEC2.SecurityGroups {
		for _, iid := range target {
			if *SecurityGroups.GroupId == iid {
				stats[iid]++                                                                  //increment hit id counter
				sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //create ec2(security group) api-instance
				_, err := sginstance.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{       //execute security group destroy
					GroupId: aws.String(iid),
				})
				if err != nil { //if there got error, print it
					fmt.Println(err)
					wg.Done()
					return
				}
				fmt.Printf("Success!\n")
			}
		}
	}

	wg.Done()
	return

}
