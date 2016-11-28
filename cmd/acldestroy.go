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

var acldestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy ACL (CAN NOT RESTORE)",
	Long: `destroy ACL (CAN NOT RESTORE)

Example:
  cq acl destroy --groupid sg-fd8cc1ee
`,
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup                //parallel processing counter group
		var gid = []string{listFlag.GroupId} //group id for printHitId
		stats := map[string]int{}            //group id hit check map
		stats[listFlag.GroupId] = 0

		if checkGroupId() != "" { //flag check
			fmt.Println(checkGroupId()) //if there is wrong, exit
			return
		}

		if listFlag.Force { //if there is enabled force option, dont confirmation
			regionsAWS := getAWSRegions() //get region list (AWS)
			for _, region := range regionsAWS {
				wg.Add(1) //waiting group count up
				go destroySecurityGroup(region, &wg, stats)
				time.Sleep(1 * time.Millisecond)
			}
			wg.Wait() //wait for end of parallel processing
		} else {
			input := ""                                                                                                 //keyboard input value
			fmt.Printf("SecurityGroup   %s   will be DESTROY, are you sure? (CAN NOT RESTORE) Y/N\n", listFlag.GroupId) //destroy warning
			fmt.Scanln(&input)                                                                                          //stdin
			if (input == "Y") || (input == "y") {                                                                       //input Y or y
				regionsAWS := getAWSRegions() //get region list (AWS)
				for _, region := range regionsAWS {
					wg.Add(1) //waiting group count up
					go destroySecurityGroup(region, &wg, stats)
					time.Sleep(1 * time.Millisecond)
				}
				wg.Wait() //wait for end of parallel processing
			} else { //not Y or y, exit
				fmt.Printf("Cancelled\n")
				return
			}
		}

		printHitId(gid, stats)

	},
}

func init() {
	aclCmd.AddCommand(acldestroyCmd)
	acldestroyCmd.Flags().StringVarP(&listFlag.GroupId, "groupid", "", "", "security group-id") // define --groupid flag
	acldestroyCmd.Flags().BoolVarP(&listFlag.Force, "force", "f", false, "Destroy without confirmation") //define -f --force flag
}

func destroySecurityGroup(region string, wg *sync.WaitGroup, stats map[string]int) {

	sgParamEC2 := getSecurityGroupParam(region) //get security group parameter

	for _, SecurityGroups := range sgParamEC2.SecurityGroups {
		if *SecurityGroups.GroupId == listFlag.GroupId {
			stats[listFlag.GroupId]++                                                     //increment hit id counter
			sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //create ec2(security group) api-instance
			_, err := sginstance.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{       //execute security group destroy
				GroupId: aws.String(listFlag.GroupId),
			})
			if err != nil { //if there got error, print it
				fmt.Println(err)
				wg.Done()
				return
			}
			fmt.Printf("Success!\n")
		}
	}

	wg.Done()
	return

}
