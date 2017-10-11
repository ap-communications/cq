package aws

import (
	"fmt"
	"sync"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Destroy(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			destroy(region, args)
		}(region)
	}
	wg.Wait()
}

func destroy(region string, target []string) {
	sgParam, err := getSecurityGroupParam(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, SecurityGroups := range sgParam.SecurityGroups {
		for _, targetId := range target {
			if *SecurityGroups.GroupId == targetId {
				if commons.Confirm() {
					sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
					if _, err := sginstance.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
						GroupId: aws.String(targetId),
					}); err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("%s was destroyed\n", targetId)
				}
			}
		}
	}
}
