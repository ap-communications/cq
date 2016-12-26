package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var dbsnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Take DB snapshot",
	Long:  "Take DB snapshot",
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

		regionsAWS := getAWSRegions()
		for _, region := range regionsAWS {
			wg.Add(1)
			go takeDBSnapshot(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {
	dbCmd.AddCommand(dbsnapshotCmd)
}

func takeDBSnapshot(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	instanceParamRDS := getRDSParam(region)
	oclock := time.Now().Format("2006-01-02-15-04-05")

	for _, DBInstances := range instanceParamRDS.DBInstances {
		for _, dbid := range target {
			if *DBInstances.DBInstanceIdentifier == dbid {
				stats[dbid]++ //increment id hit counter
				rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
				_, err := rdsInstance.CreateDBSnapshot(&rds.CreateDBSnapshotInput{
					DBInstanceIdentifier: aws.String(dbid),
					DBSnapshotIdentifier: aws.String("cq-" + oclock),
				})
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("Success!   %s   has started take snapshot\n", dbid)
			}
		}
	}

	return
}
