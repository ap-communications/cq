package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"strings"
	"sync"
	"time"
)

var dbdestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "[DANGER]  Destroy DB  (CAN NOT RESTORE)",
	Long:  "[DANGER]  Destroy DB  (CAN NOT RESTORE)",
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
			//jump to destroy sequence
		} else {
			input := ""                                                        //keyboard input value
			fmt.Printf("DB   %s   will be destroy, are you sure?  Y/N\n", ids) //destroy warning (pre)
			fmt.Scanln(&input)                                                 //stdin
			if (input == "Y") || (input == "y") {                              //input Y or y
				fmt.Printf("This is final warning. DESTROY DB   %s   ARE YOU SURE?  Y/N\n", ids) //final destroy warning
				fmt.Scanln(&input)                                                               //stdin
				if (input == "Y") || (input == "y") {                                            //input Y or y (2 (final))
					//jump to destroy sequence
				} else {
					fmt.Printf("Cancelled\n")
					return
				}
			} else {
				fmt.Printf("Cancelled\n")
				return
			}
		}

		var wg sync.WaitGroup
		regionsAWS := getAWSRegions()
		for _, region := range regionsAWS {
			wg.Add(1)
			go destroyDB(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {
	dbCmd.AddCommand(dbdestroyCmd)
	dbdestroyCmd.Flags().BoolVarP(&listFlag.Force, "force", "f", false, "Destroy without confirmation") //define -f --force flag
}

func destroyDB(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	instanceParamRDS := getRDSParam(region)

	for _, DBInstances := range instanceParamRDS.DBInstances {
		for _, dbid := range target {
			if *DBInstances.DBInstanceIdentifier == dbid {
				stats[dbid]++
				rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
				_, err := rdsInstance.DeleteDBInstance(&rds.DeleteDBInstanceInput{
					DBInstanceIdentifier: aws.String(dbid),
					SkipFinalSnapshot:    aws.Bool(true),
				})
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("Success!   %s   destroyed\n", dbid)
			}
		}
	}

	return

}
