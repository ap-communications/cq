package aws

import (
	"fmt"
	"strings"
	"sync"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func List(ch chan<- commons.InstanceList) {
	defer close(ch)

	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			i, err := getEC2Instances(region)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			sendInstanceList(ch, i)
		}(region)
	}
	wg.Wait()
	return
}

func sendInstanceList(ch chan<- commons.InstanceList, instance ec2.DescribeInstancesOutput) {
	for _, Reservations := range instance.Reservations {
		for _, Instances := range Reservations.Instances {
			var l commons.InstanceList

			if Instances.InstanceId == nil {
				l.InstanceId = "NULL"
			} else {
				l.InstanceId = *Instances.InstanceId
			}

			if Instances.State.Name == nil {
				l.State = "NULL"
			} else {
				l.State = *Instances.State.Name
			}

			if Instances.PublicIpAddress == nil {
				l.PublicIpAddress = "NULL"
			} else {
				l.PublicIpAddress = *Instances.PublicIpAddress
			}

			if Instances.PrivateIpAddress == nil {
				l.PrivateIpAddress = "NULL"
			} else {
				l.PrivateIpAddress = *Instances.PrivateIpAddress
			}

			if Instances.Placement.AvailabilityZone == nil {
				l.AvailabilityZone = "NULL"
			} else {
				l.AvailabilityZone = *Instances.Placement.AvailabilityZone
			}

			if Instances.Tags == nil {
				l.Tags = "NULL"
			} else if strings.Contains(*Instances.Tags[0].Value, commons.Flags.Delimiter) && (commons.Flags.Delimiter != "") { //When the delimiter is included in Name-Tag, enclose it with double quotation because separation will increase
				l.Tags = "\"" + *Instances.Tags[0].Value + "\"" //enclose Name-Tag with double quotation
			} else {
				l.Tags = *Instances.Tags[0].Value
			}

			l.Provider = "AWS"

			ch <- l
		}
	}
	return
}
