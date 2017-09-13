package aws

import (
	"fmt"
	"sync"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func StopInstance(args []string) {
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			stopInstance(args, region)
		}(region)
	}
	wg.Wait()
}

func stopInstance(targets []string, region string) {
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
						resp, err := ec2instance.StopInstances(&ec2.StopInstancesInput{InstanceIds: []*string{aws.String(target)}})
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						fmt.Printf("%s %s ===> %s\n", target, *resp.StoppingInstances[0].PreviousState.Name, *resp.StoppingInstances[0].CurrentState.Name)
					}
				}
			}
		}
	}
}
