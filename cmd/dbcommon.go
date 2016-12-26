package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"sync"
)

func setRDSParam(region string, wg *sync.WaitGroup, instanceParamRDS *[]rds.DescribeDBInstancesOutput) {

	defer wg.Done()

	rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := rdsInstance.DescribeDBInstances(&rds.DescribeDBInstancesInput{})

	if err != nil {
		fmt.Println(err)
		return
	}

	*instanceParamRDS = append(*instanceParamRDS, *resp) //set response instanceParamEC2 array

	return

}

func getRDSParam(region string) *rds.DescribeDBInstancesOutput {

	rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := rdsInstance.DescribeDBInstances(&rds.DescribeDBInstancesInput{})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp
}

func setRDSSnapList(region string, wg *sync.WaitGroup, snapshotParamRDS *[]rds.DescribeDBSnapshotsOutput) {

	defer wg.Done()

	rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := rdsInstance.DescribeDBSnapshots(&rds.DescribeDBSnapshotsInput{})

	if err != nil {
		fmt.Println(err)
		return
	}

	*snapshotParamRDS = append(*snapshotParamRDS, *resp) //set response instanceParamEC2 array

	return

}
