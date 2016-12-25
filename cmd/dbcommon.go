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
