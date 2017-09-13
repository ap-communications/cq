package aws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func RestoreInstance(args []string) {
	var wg sync.WaitGroup
	for _, filepath := range args {
		wg.Add(1)
		go func(filepath string) {
			defer wg.Done()
			restoreInstance(filepath)
		}(filepath)
	}
	wg.Wait()
}

func restoreInstance(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data ec2.Instance
	unmarshalErr := json.Unmarshal([]byte(string(file)), &data)
	if unmarshalErr != nil {
		fmt.Println(err)
		return
	}

	region := *data.Placement.AvailabilityZone
	azSize := len(region) - 1 //AZ format is region name + one alphabet
	region = region[:azSize]

	monitoringEnabled := true
	if *data.Monitoring.State == "disabled" {
		monitoringEnabled = false
	}

	var iamProfile string
	if data.IamInstanceProfile != nil {
		iamProfile = *data.IamInstanceProfile.Arn
	}

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := ec2instance.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(*data.ImageId),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		EbsOptimized: aws.Bool(*data.EbsOptimized),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(iamProfile),
		},
		InstanceType: aws.String(*data.InstanceType),
		KeyName:      aws.String(*data.KeyName),
		Monitoring: &ec2.RunInstancesMonitoringEnabled{
			Enabled: aws.Bool(monitoringEnabled),
		},
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(*data.Placement.AvailabilityZone),
			Tenancy:          aws.String(*data.Placement.Tenancy),
		},
		SecurityGroupIds: []*string{
			aws.String(*data.SecurityGroups[0].GroupId),
		},
		SubnetId: aws.String(*data.SubnetId),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	checkedResp := checkEC2InstanceCreated(region, *resp.Instances[0].InstanceId)
	instance := checkedResp.Reservations[0].Instances[0]
	fmt.Printf("%s %s restored\n", *instance.InstanceId, *instance.PublicIpAddress)
}
