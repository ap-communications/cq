package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getEC2Instances(region string) (ec2.DescribeInstancesOutput, error) {
	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := ec2instance.DescribeInstances(nil)
	if err != nil {
		fmt.Println(err)
		return *resp, err
	}
	return *resp, nil
}

func createKey(region string) (string, string) {
	oclock := time.Now().Format("2006-01-02_15:04:05")
	name := "cq_" + oclock
	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := ec2instance.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(name),
	})
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	return *resp.KeyMaterial, name
}

func checkEC2InstanceCreated(region string, id string) *ec2.DescribeInstancesOutput {
	resp := getEC2Instance(region, id)
	for {
		for _, Reservations := range resp.Reservations {
			for _, Instances := range Reservations.Instances {
				if Instances.PublicIpAddress != nil {
					fmt.Printf("\n")
					return resp
				}
			}
		}
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf(".")
		resp = getEC2Instance(region, id)
	}
}

func getEC2Instance(region string, id string) *ec2.DescribeInstancesOutput {
	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := ec2instance.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp
}
