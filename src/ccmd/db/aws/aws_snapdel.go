package aws

import (
	"fmt"
	"sync"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func SnapDel(args []string) {
	if len(args) == 0 { //If there is no argument, abort
		fmt.Printf("missing args (DB-ID)\n")
		return
	}
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			snapDel(region, args)
		}(region)
	}
	wg.Wait()
}

func snapDel(region string, target []string) {
	snapList, err := getRdsSnapList(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, DBSnapshots := range snapList.DBSnapshots {
		for _, snapId := range target {
			if *DBSnapshots.DBSnapshotIdentifier == snapId {
				if commons.Confirm() {
					rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
					_, err := rdsInstance.DeleteDBSnapshot(&rds.DeleteDBSnapshotInput{
						DBSnapshotIdentifier: aws.String(snapId),
					})
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("%s was deleted\n", snapId)
				}
			}
		}
	}
}
