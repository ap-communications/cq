package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sync"
	"time"
)

var dbbackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "DB Backup configure information",
	Long:  "DB Backup configure information",
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

		regionsAWS := getAWSRegions() //get region list
		for _, region := range regionsAWS {
			wg.Add(1) //waiting group count up
			go backupDBInstance(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {

	USAGE := `Usage:
  cq db backup [instance-id] [instance-id] ...
`

	dbCmd.AddCommand(dbbackupCmd)
	dbbackupCmd.SetUsageTemplate(USAGE) //override Usage words
}

func backupDBInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	instanceParamRDS := getRDSParam(region)

	for _, DBInstances := range instanceParamRDS.DBInstances {
		for _, iid := range target {
			if *DBInstances.DBInstanceIdentifier == iid {
				stats[iid]++ //increment id hit counter
				file, err := os.Create("./" + *DBInstances.DBInstanceIdentifier + ".aws")
				defer file.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
				data, _ := json.Marshal(DBInstances)
				file.Write([]byte(string(data)))
				fmt.Printf("backuped  %s\n", *DBInstances.DBInstanceIdentifier+".aws")
			}
		}
	}

	return

}
