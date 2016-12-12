package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strconv"
	"sync"
	"time"
)

func setSecurityGroupParam(region string, wg *sync.WaitGroup, sgParamEC2 *[]ec2.DescribeSecurityGroupsOutput) {

	defer wg.Done()

	sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := sginstance.DescribeSecurityGroups(nil)                           //get instance information from API

	//エラーがあれば出力し終了
	if err != nil {
		fmt.Println(err)
		return
	}

	*sgParamEC2 = append(*sgParamEC2, *resp) //set response data sgParamEC2 array

	return

}

func getSecurityGroupParam(region string) *ec2.DescribeSecurityGroupsOutput {

	sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := sginstance.DescribeSecurityGroups(nil)                           //get instance information from API

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp //return response data (json)
}

func checkSecurityGroupFlags() string {

	checkResult := ""

	checkResult += checkGroupId()
	checkResult += checkProtocol()
	checkResult += checkWay()

	return checkResult

}

func checkGroupId() string {

	if listFlag.GroupId == "" { //check security group id is not empty
		return "required group-id   --groupid string\n"
	}

	return ""

}

func checkProtocol() string {

	if (listFlag.Protocol != "tcp") && (listFlag.Protocol != "udp") && (listFlag.Protocol != "icmp") && (listFlag.Protocol != "any") && (listFlag.Protocol != "") { //check protocol is not empty, or unexpected
		return "invalid protocol (tcp, udp, icmp, any)\n"
	}

	return ""

}

func checkWay() string {

	if listFlag.Way == "" { //check way of packet is not empy, or invalid
		return "required way of packet   --way string\n"
	} else if (listFlag.Way != "ingress") && (listFlag.Way != "egress") {
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

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
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
