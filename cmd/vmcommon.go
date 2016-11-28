package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sync"
	"time"
)

func setEC2Param(region string, wg *sync.WaitGroup, instanceParamEC2 *[]ec2.DescribeInstancesOutput) {

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := ec2instance.DescribeInstances(nil)                                //get instance information from API

	if err != nil {
		fmt.Println(err)
		return
	}

	*instanceParamEC2 = append(*instanceParamEC2, *resp) //set response instanceParamEC2 array

	wg.Done()

}

func getEC2Param(region string) *ec2.DescribeInstancesOutput {

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := ec2instance.DescribeInstances(nil)                                //get instance information from API

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp //return response (json)
}

func createKey(region string) (string, string) {

	name := "cq_" + time.Now().String()

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := ec2instance.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(name),
	})

	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	return *resp.KeyMaterial, name //return secret-key

}

func getInstanceParam(region string, iid string) *ec2.DescribeInstancesOutput {

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := ec2instance.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(iid),
		},
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp //return response (json)
}

func checkInstanceCreated(region string, iid string) *ec2.DescribeInstancesOutput {

	resp := getInstanceParam(region, iid)

	for i := 0; i < 60; i++ { //60 seconds

		for _, Reservations := range resp.Reservations {
			for _, Instances := range Reservations.Instances {
				if Instances.PublicIpAddress != nil {
					fmt.Printf("\n")
					return resp
				}
			}
		}

		time.Sleep(1000 * time.Millisecond)
		resp = getInstanceParam(region, iid)

	}

	return resp
}
