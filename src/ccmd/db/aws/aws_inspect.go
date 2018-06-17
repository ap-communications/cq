package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/service/rds"
)

func Inspect(args []string) {
	regions := commons.GetAwsRegions()
	m := map[string]rds.DescribeDBInstancesOutput{}
	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			i, err := getRdsInstances(region)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			m[region] = i
		}(region)
	}
	wg.Wait()

	var rdsInstances []rds.DescribeDBInstancesOutput
	for _, region := range regions {
		rdsInstances = append(rdsInstances, m[region])
	}
	printDbInspect(rdsInstances, args)
}
func printDbInspect(rdsInstances []rds.DescribeDBInstancesOutput, args []string) {
	if len(args) == 0 { //If there is no argument, print all DB instances information
		for _, lawJson := range rdsInstances {
			for _, DBInstances := range lawJson.DBInstances {
				fmt.Println(DBInstances)
			}
		}
	} else {
		for _, lawJson := range rdsInstances {
			for _, DBInstances := range lawJson.DBInstances {
				for _, dbId := range args {
					if *DBInstances.DBInstanceIdentifier == dbId {
						fmt.Println(*DBInstances)
					}
				}
			}
		}
	}
}
