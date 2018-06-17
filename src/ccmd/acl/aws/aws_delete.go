package aws

import (
	"fmt"
	"sync"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Delete(args []string) {
	if checkSecurityGroupFlags() != "" {
		fmt.Println(checkSecurityGroupFlags())
		return
	}
	var wg sync.WaitGroup
	for _, region := range commons.GetAwsRegions() {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			delete(region)
		}(region)
	}
	wg.Wait()
}

func delete(region string) {
	sgParamEC2, err := getSecurityGroupParam(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	port := translateSecurityGroupPort(commons.Flags.Port)
	protocol := translateSecurityGroupProtocol(commons.Flags.Protocol)
	address := translateSecurityGroupAddress(commons.Flags.Address)

	for _, SecurityGroups := range sgParamEC2.SecurityGroups {
		if *SecurityGroups.GroupId == commons.Flags.GroupId {
			sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
			switch {
			case commons.Flags.Way == "ingress":
				if _, err := sginstance.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
					GroupId: aws.String(commons.Flags.GroupId),
					IpPermissions: []*ec2.IpPermission{
						{
							FromPort:   aws.Int64(port),
							ToPort:     aws.Int64(port),
							IpProtocol: aws.String(protocol),
							IpRanges: []*ec2.IpRange{
								{CidrIp: aws.String(address)},
							},
						},
					},
				}); err != nil {
					fmt.Println(err)
					return
				}
			case commons.Flags.Way == "egress":
				if _, err := sginstance.RevokeSecurityGroupEgress(&ec2.RevokeSecurityGroupEgressInput{
					GroupId: aws.String(commons.Flags.GroupId),
					IpPermissions: []*ec2.IpPermission{
						{
							FromPort:   aws.Int64(port),
							ToPort:     aws.Int64(port),
							IpProtocol: aws.String(protocol),
							IpRanges: []*ec2.IpRange{
								{CidrIp: aws.String(address)},
							},
						},
					},
				}); err != nil {
					fmt.Println(err)
					return
				}
			}
			fmt.Printf("%s Done\n", commons.Flags.GroupId)
		}
	}
}
