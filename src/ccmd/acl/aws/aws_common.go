package aws

import (
	"fmt"
	"strconv"
	"time"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSecurityGroupParam(region string) (ec2.DescribeSecurityGroupsOutput, error) {
	sgInstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := sgInstance.DescribeSecurityGroups(nil)
	if err != nil {
		return ec2.DescribeSecurityGroupsOutput{}, err
	}
	return *resp, nil
}

//check setting security group flags is not unexpected
func checkSecurityGroupFlags() string {
	checkResult := ""
	checkResult += checkGroupId()
	checkResult += checkProtocol()
	checkResult += checkWay()
	return checkResult
}

//check security group id is not empty
func checkGroupId() string {
	if commons.Flags.GroupId == "" {
		return "required group-id   --groupid string\n"
	}
	return ""
}

//check protocol is not empty, or unexpected
func checkProtocol() string {
	if (commons.Flags.Protocol != "tcp") && (commons.Flags.Protocol != "udp") && (commons.Flags.Protocol != "icmp") && (commons.Flags.Protocol != "any") && (commons.Flags.Protocol != "") {
		return "invalid protocol (tcp, udp, icmp, any)\n"
	}
	return ""
}

//check way of packet is not empy, or invalid
func checkWay() string {
	if commons.Flags.Way == "" {
		return "required way of packet   --way string\n"
	} else if (commons.Flags.Way != "ingress") && (commons.Flags.Way != "egress") {
		return "invalid way of packet (ingress or egress)\n"
	}
	return ""
}

func translateSecurityGroupPort(argPort string) int64 {
	var port int64
	if (argPort == "any") || (argPort == "") { // -1 is "any" in AWS
		port = -1
	} else {
		port, _ = strconv.ParseInt(argPort, 10, 64)
	}
	return port
}

func translateSecurityGroupProtocol(argProtocol string) string {
	var protocol string
	if (argProtocol == "any") || (argProtocol == "") { // -1 is "any" in AWS
		protocol = "-1"
	} else {
		protocol = argProtocol
	}
	return protocol
}

func translateSecurityGroupAddress(argAddress string) string {
	var address string
	if (argAddress == "any") || (argAddress == "") { //"any" or not set default, 0.0.0.0/0
		address = "0.0.0.0/0"
	} else {
		address = argAddress
	}
	return address
}

func createSecurityGroup(region string) string {
	oclock := time.Now().Format("2006-01-02_15:04:05")
	groupname := "cq_temporary_sg_" + oclock
	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := ec2instance.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description: aws.String("Created by cq    " + oclock),
		GroupName:   aws.String(groupname),
	})
	if err != nil {
		fmt.Println(err)
		return "err"
	}
	return *resp.GroupId
}
