package aws

import (
	"fmt"
	"sync"
	"time"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func Snapshot(args []string) {
	if len(args) == 0 { //If there is no argument, abort
		fmt.Printf("missing args (DB-ID)\n")
		return
	}
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			snapshot(region, args)
		}(region)
	}
	wg.Wait()
}

func snapshot(region string, target []string) {
	rdsInstances, err := getRdsInstances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	now := time.Now().Format("2006-01-02-15-04-05")
	for _, DBInstances := range rdsInstances.DBInstances {
		for _, dbId := range target {
			if *DBInstances.DBInstanceIdentifier == dbId {
				rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
				if _, err := rdsInstance.CreateDBSnapshot(&rds.CreateDBSnapshotInput{
					DBInstanceIdentifier: aws.String(dbId),
					DBSnapshotIdentifier: aws.String("cq-" + now),
				}); err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("%s  has started take snapshot\n", dbId)
			}
		}
	}
}
