package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func Reboot(args []string) {
	if len(args) == 0 { //If there is no argument, abort
		fmt.Printf("missing args (DB-ID)\n")
		return
	}
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			reboot(region, args)
		}(region)
	}
	wg.Wait()
}

func reboot(region string, target []string) {
	rdsInstances, err := getRdsInstances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, DBInstances := range rdsInstances.DBInstances {
		for _, dbId := range target {
			if *DBInstances.DBInstanceIdentifier == dbId {
				if *DBInstances.DBInstanceStatus == "available" { //if instance state is not "available" can't be reboot)
					if commons.Confirm() {
						rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
						if _, err := rdsInstance.RebootDBInstance(&rds.RebootDBInstanceInput{
							DBInstanceIdentifier: aws.String(dbId),
							ForceFailover:        aws.Bool(commons.Flags.NoFailover),
						}); err != nil {
							fmt.Println(err)
							return
						}
						fmt.Printf("%s has started reboot sequence\n", dbId)
					}
				} else {
					fmt.Printf("%s Can't reboot, DB instance status is now %s\n", dbId, *DBInstances.DBInstanceStatus)
				}
			}
		}
	}
}
