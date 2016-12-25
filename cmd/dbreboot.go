package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"sync"
	"strings"
	"time"
)

var dbrebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot DB with failover (default)",
	Long:  `Reboot DB with failover (default)
If you want to no failover reboot, must set failover flag at false


*** WARNING ***
  If you execute to no failover reboot, it will be take more DB down time


INFO:
  If you want not configured MultiAZ or High Availability to DB instance, you must set --no-failover flag



...
`,
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup     //parallel processing counter group
		stats := map[string]int{} //instance id hit check map

		for _, argid := range args { //create hit judgment map for character string set as argument
			stats[argid] = 0 //init map (hit is 0)
		}

		if len(args) == 0 { //If there is no argument, abort
			fmt.Printf("missing args (DB-ID)\n")
			return
		}

		ids := ""                      //keyboard input value
		for _, inputid := range args { //translate comma spreaded (for warning print)
			ids += inputid + ", "
		}
		ids = strings.TrimRight(ids, ", ") //delete final comma

		if listFlag.Force { //if there is enabled force option, dont confirmation
			//jump to reboot sequence
		} else {
			input := ""                                                              //keyboard input value
			fmt.Printf("DB   %s   will be reboot, are you sure?  Y/N\n", ids) //destroy warning (pre)
			fmt.Scanln(&input)                                                       //stdin
			if (input == "Y") || (input == "y") {                                    //input Y or y
				//jump to reboot sequence
			} else {
				fmt.Printf("Cancelled\n")
				return
			}
		}

		regionsAWS := getAWSRegions()
		for _, region := range regionsAWS {
			wg.Add(1)
			go rebootDBInstance(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {
	dbCmd.AddCommand(dbrebootCmd)
	dbrebootCmd.Flags().BoolVarP(&dbFlag.Failover, "no-failover", "", false, "execute to no failover reboot (it will be take more DB down time)")                    // define --failover flag
}

func rebootDBInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	var failover bool
	defer wg.Done()

	instanceParamRDS := getRDSParam(region)

	/*
	set no-failover flag status for AWS API (RebootDBInstance method is can set ForceFailover only. must be reverse dbFlag.Failover)
		* RebootDBInstance method option
			ForceFailover:  true = no-failover,   false = with-failover
		* dbFlag.Failover
			flag set: true,   not set: false
	*/
	if dbFlag.Failover {
		failover = false
	} else {
		failover = true
	}

	for _, DBInstances := range instanceParamRDS.DBInstances {
		for _, dbid := range target {
			if *DBInstances.DBInstanceIdentifier == dbid {
				stats[dbid]++ //increment id hit counter
				if *DBInstances.DBInstanceStatus == "available" { //if RDS instance state is "available" (coz: if instance state is not "available" can't be reboot)
					rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
					_, err := rdsInstance.RebootDBInstance(&rds.RebootDBInstanceInput{
						DBInstanceIdentifier: aws.String(dbid),
						ForceFailover: aws.Bool(failover),
					})
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("Success!   %s   has started reboot sequence\n", dbid)
				} else {
					fmt.Printf("Can't reboot, DB instance   %s   status is now %s\n", dbid, *DBInstances.DBInstanceStatus) //if RDS instance state is not "available" can't be reboot
				}
			}
		}
	}

	return
}
