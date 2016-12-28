package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sync"
	"time"
)

var vmbackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup configure information",
	Long:  "Backup configure information",
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
			go backupInstance(args, region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		printHitId(args, stats)

	},
}

func init() {

	USAGE := `Usage:
  cq vm backup [instance-id] [instance-id] ...
`

	vmCmd.AddCommand(vmbackupCmd)
	vmstartCmd.SetUsageTemplate(USAGE) //override Usage words
}

func backupInstance(target []string, region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	instanceParamEC2 := getEC2Param(region)

	for _, Reservations := range instanceParamEC2.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, iid := range target {
				if *Instances.InstanceId == iid {
					var nameTag string
					if len(Instances.Tags) == 0 {
						nameTag = ""
					} else {
						tagmap := map[string]string{}
						for _, Tags := range Instances.Tags {
							tagmap[*Tags.Key] = *Tags.Value
						}
						nameTag = tagmap["Name"]
					}
					file, err := os.Create("./" + *Instances.InstanceId + "_" + nameTag + ".aws")
					if err != nil {
						fmt.Println(err)
					}
					defer file.Close()
					data, _ := json.Marshal(Instances)
					file.Write([]byte(string(data)))
					fmt.Printf("backuped  %s\n", *Instances.InstanceId+"_"+nameTag+".aws")
					stats[iid]++ //increment id hit counter
				}
			}
		}
	}

	return

}
