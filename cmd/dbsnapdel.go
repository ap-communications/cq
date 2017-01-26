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

var dbsnapdelCmd = &cobra.Command{
	Use:   "snapdel",
	Short: "Delete DB snapshot",
	Long:  "Delete DB snapshot",
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

		if listFlag.Force { //if there is enabled force option, dont confirmation
			//jump to delete sequence
		} else {
			input := ""                                                              //keyboard input value
			fmt.Printf("Snapshot   %s   will be delete, are you sure?  Y/N\n", args) //destroy warning (pre)
			fmt.Scanln(&input)                                                       //stdin
			if (input == "Y") || (input == "y") {                                    //input Y or y
				//jump to delete sequence
			} else {
				fmt.Printf("Cancelled\n")
				return
			}
		}

		regionsAWS := getAWSRegions()
		for _, region := range regionsAWS {
			wg.Add(1)
			go deleteDBSnapshot(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {
	dbCmd.AddCommand(dbsnapdelCmd)
	dbsnapdelCmd.Flags().BoolVarP(&listFlag.Force, "force", "f", false, "Delete without confirmation") //define -f --force flag
}

func deleteDBSnapshot(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	snapshotParamRDS := getRDSSnapList(region)

	for _, DBSnapshots := range snapshotParamRDS.DBSnapshots {
		for _, snapid := range target {
			if *DBSnapshots.DBSnapshotIdentifier == snapid {
				stats[snapid]++ //increment id hit counter
				rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
				_, err := rdsInstance.DeleteDBSnapshot(&rds.DeleteDBSnapshotInput{
					DBSnapshotIdentifier: aws.String(snapid),
				})
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("Success!   %s   was deleted\n", snapid)
			}
		}
	}

	return

}
