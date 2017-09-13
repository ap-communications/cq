package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"ccmd/commons"
)

func BackupInstance(args []string) {
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			backupInstance(args, region)
		}(region)
	}
	wg.Wait()
}

func backupInstance(targets []string, region string) {
	instances, err := getEC2Instances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, Reservations := range instances.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, target := range targets {
				if *Instances.InstanceId == target {
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
					filePath := *Instances.InstanceId + "_" + nameTag + ".cq"
					file, err := os.Create("./" + filePath)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer file.Close()
					data, _ := json.Marshal(Instances)
					file.Write([]byte(string(data)))
					fmt.Printf("%s\n", filePath)
				}
			}
		}
	}
}
