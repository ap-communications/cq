package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func getRdsInstances(region string) (rds.DescribeDBInstancesOutput, error) {
	rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := rdsInstance.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
	if err != nil {
		return rds.DescribeDBInstancesOutput{}, err
	}
	return *resp, nil
}

func getRdsSnapList(region string) (rds.DescribeDBSnapshotsOutput, error) {
	rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := rdsInstance.DescribeDBSnapshots(&rds.DescribeDBSnapshotsInput{})
	if err != nil {
		return rds.DescribeDBSnapshotsOutput{}, err
	}
	return *resp, nil
}
