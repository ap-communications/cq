package aws

import (
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type instanceList struct {
	Tags             string
	InstanceId       string
	State            string
	PublicIpAddress  string
	PrivateIpAddress string
	AvailabilityZone string
	Provider         string
}

func PrintVmList(w *tabwriter.Writer, column string) {
	m := map[string]ec2.DescribeInstancesOutput{}
	regions := commons.GetAwsRegions()

	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			i, err := getEC2Instances(region)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			m[region] = i
		}(region)
	}
	wg.Wait()

	var instances []ec2.DescribeInstancesOutput
	for _, region := range regions {
		instances = append(instances, m[region])
	}
	instanceLists := getInstanceList(instances)
	inject(w, column, instanceLists)
}

func getInstanceList(instances []ec2.DescribeInstancesOutput) []instanceList {
	var instanceLists []instanceList
	for _, lawJson := range instances {
		for _, Reservations := range lawJson.Reservations {
			for _, Instances := range Reservations.Instances {
				var l instanceList

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

				instanceLists = append(instanceLists, l)
			}
		}
	}
	return instanceLists
}

func inject(w *tabwriter.Writer, column string, instanceLists []instanceList) {
	for _, l := range instanceLists {
		fmt.Fprintf(
			w,
			column,
			l.Tags,
			l.InstanceId,
			l.State,
			l.PublicIpAddress,
			l.PrivateIpAddress,
			l.AvailabilityZone,
			l.Provider,
		)
	}
}
