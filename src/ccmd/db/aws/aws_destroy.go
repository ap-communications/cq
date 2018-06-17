package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func Destroy(args []string) {
	if len(args) == 0 { //If there is no argument, abort
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			destroy(region, args)
		}(region)
	}
	wg.Wait()
}

func destroy(region string, target []string) {
	rdsInstances, err := getRdsInstances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, DBInstances := range rdsInstances.DBInstances {
		for _, dbId := range target {
			if *DBInstances.DBInstanceIdentifier == dbId {
				if commons.Confirm() {
					rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
					if _, err := rdsInstance.DeleteDBInstance(&rds.DeleteDBInstanceInput{
						DBInstanceIdentifier: aws.String(dbId),
						SkipFinalSnapshot:    aws.Bool(true),
					}); err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("%s was destroyed\n", dbId)
				}
			}
		}
	}
}
