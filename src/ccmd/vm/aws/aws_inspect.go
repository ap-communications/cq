package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"
)

func Inspect(args []string) {
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			inspect(args, region)
		}(region)
	}
	wg.Wait()
}

func inspect(targets []string, region string) {
	instances, err := getEC2Instances(region)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, Reservations := range instances.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, target := range targets {
				if *Instances.InstanceId == target {
					fmt.Println(Instances)
				}
			}
		}
	}
}
