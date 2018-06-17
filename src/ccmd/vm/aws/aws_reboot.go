package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func RebootInstance(args []string) {
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			rebootInstance(args, region)
		}(region)
	}
	wg.Wait()
}

func rebootInstance(targets []string, region string) {
	instances, err := getEC2Instances(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, Reservations := range instances.Reservations {
		for _, Instances := range Reservations.Instances {
			for _, target := range targets {
				if *Instances.InstanceId == target {
					if *Instances.State.Name != "running" {
						fmt.Printf("Can't be reboot %s %s\n", target, *Instances.State.Name)
						return
					}
					if commons.Confirm() {
						ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
						_, err := ec2instance.RebootInstances(&ec2.RebootInstancesInput{InstanceIds: []*string{aws.String(target)}})
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						fmt.Printf("%s has started reboot sequence\n", target)
					}
				}
			}
		}
	}
}
