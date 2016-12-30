package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var dbinspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Print detailed DB instance info",
	Long:  "Print detailed DB instance info",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup
		var instanceParamRDS = []rds.DescribeDBInstancesOutput{}

		stats := map[string]int{}    //instance id hit check map
		for _, argid := range args { //create hit judgment map for character string set as argument
			stats[argid] = 0 //init map (hit count is 0)
		}

		regionsAWS := getAWSRegions() //get region list (AWS)
		for _, region := range regionsAWS {
			wg.Add(1)
			go setRDSParam(region, &wg, &instanceParamRDS)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printDBInspect(instanceParamRDS, args, stats)
		printHitId(args, stats)

	},
}

func init() {

	USAGE := `Usage:
  cq db inspect
    or
  cq db inspect [DB-id] [DB-id] ...
`

	dbCmd.AddCommand(dbinspectCmd)
	dbinspectCmd.SetUsageTemplate(USAGE) //override Usage words
}

func printDBInspect(instanceParamRDS []rds.DescribeDBInstancesOutput, args []string, stats map[string]int) {

	if len(args) == 0 { //If there is no argument, print all DBinstances information
		for _, lawJson := range instanceParamRDS {
			for _, DBInstances := range lawJson.DBInstances {
				fmt.Println(DBInstances)
			}
		}
	} else { //if there is argument DBinstance id, print it
		for _, lawJson := range instanceParamRDS {
			for _, DBInstances := range lawJson.DBInstances {
				for _, dbid := range args {
					if *DBInstances.DBInstanceIdentifier == dbid {
						fmt.Println(*DBInstances)
						stats[dbid]++ //increment hit id counter
					}
				}
			}
		}
	}

	return

}
