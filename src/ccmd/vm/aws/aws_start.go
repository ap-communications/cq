package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func StartInstance(args []string) {
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			startInstance(args, region)
		}(region)
	}
	wg.Wait()
}

func startInstance(targets []string, region string) {
	instances, err := getEC2Instances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, Reservations := range instances.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, target := range targets {
				if *Instances.InstanceId == target {
					ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
					resp, err := ec2instance.StartInstances(&ec2.StartInstancesInput{InstanceIds: []*string{aws.String(target)}})
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					fmt.Printf("%s %s ===> %s\n", target, *resp.StartingInstances[0].PreviousState.Name, *resp.StartingInstances[0].CurrentState.Name)
				}
			}
		}
	}
}
