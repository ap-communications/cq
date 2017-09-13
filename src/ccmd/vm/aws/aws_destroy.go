package aws

import (
	"fmt"
	"sync"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DestroyInstance(args []string) {
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			destroyInstance(args, region)
		}(region)
	}
	wg.Wait()
}

func destroyInstance(targets []string, region string) {
	instances, err := getEC2Instances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, Reservations := range instances.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, target := range targets {
				if *Instances.InstanceId == target {
					if commons.Confirm() {
						ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
						_, err := ec2instance.TerminateInstances(&ec2.TerminateInstancesInput{InstanceIds: []*string{aws.String(target)}})
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						fmt.Printf("%s was destroyed\n", target)
					}
				}
			}
		}
	}
}
