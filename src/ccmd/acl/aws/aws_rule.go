package aws

import (
	"fmt"
	"strconv"
	"sync"
	"text/tabwriter"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type ruleList struct {
	Way      string //way of packet
	Protocol string
	Port     string
	Address  string
}

func RuleList(w *tabwriter.Writer, column string, args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (SecurityGroup-ID)\n")
		return
	}
	m := map[string]ec2.DescribeSecurityGroupsOutput{}
	regions := commons.GetAwsRegions()
	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			r, err := getSecurityGroupParam(region)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			m[region] = r
		}(region)
	}
	wg.Wait()

	var rules []ec2.DescribeSecurityGroupsOutput
	for _, region := range regions {
		rules = append(rules, m[region])
	}
	ruleLists := getSecurityGroupRules(rules, args)
	injectRuleList(w, column, ruleLists)
}

func getSecurityGroupRules(sgParams []ec2.DescribeSecurityGroupsOutput, target []string) []ruleList {
	var ruleLists []ruleList
	for _, sgParam := range sgParams {
		for _, SecurityGroups := range sgParam.SecurityGroups {
			for _, sgId := range target {
				if *SecurityGroups.GroupId == sgId {
					for _, IpPermissions := range SecurityGroups.IpPermissions {
						var l ruleList

						l.Way = "Ingress"

						if IpPermissions.IpProtocol == nil {
							l.Protocol = "NULL"
						} else if *IpPermissions.IpProtocol == "-1" {
							l.Protocol = "any"
						} else {
							l.Protocol = *IpPermissions.IpProtocol
						}

						if IpPermissions.FromPort == nil {
							l.Port = "any"
						} else if *IpPermissions.IpProtocol == "icmp" { //if protocol is icmp, format is Type & Code
							if *IpPermissions.FromPort == -1 {
								l.Port += "Type:any"
							} else {
								l.Port += "Type:" + strconv.FormatInt(*IpPermissions.FromPort, 10)
							}
							if *IpPermissions.ToPort == -1 {
								l.Port += "_Code:any"
							} else {
								l.Port += "_Code:" + strconv.FormatInt(*IpPermissions.ToPort, 10)
							}
						} else {
							l.Port = strconv.FormatInt(*IpPermissions.FromPort, 10)
						}

						if IpPermissions.IpRanges == nil {
							l.Address = "NULL"
						} else {
							for _, IpRanges := range IpPermissions.IpRanges {
								l.Address = *IpRanges.CidrIp
								ruleLists = append(ruleLists, l)
							}
						}

						if IpPermissions.UserIdGroupPairs == nil {
							l.Address = "NULL"
						} else {
							for _, UserIdGroupPairs := range IpPermissions.UserIdGroupPairs {
								l.Address = *UserIdGroupPairs.GroupId
								ruleLists = append(ruleLists, l)
							}
						}

					}

					for _, IpPermissionsEgress := range SecurityGroups.IpPermissionsEgress {
						var l ruleList

						l.Way = "Egress"

						if IpPermissionsEgress.IpProtocol == nil {
							l.Protocol = "NULL"
						} else if *IpPermissionsEgress.IpProtocol == "-1" {
							l.Protocol = "any"
						} else {
							l.Protocol = *IpPermissionsEgress.IpProtocol
						}

						if IpPermissionsEgress.FromPort == nil {
							l.Port = "any"
						} else if *IpPermissionsEgress.IpProtocol == "icmp" { //if protocol is icmp, format is Type & Code
							if *IpPermissionsEgress.FromPort == -1 {
								l.Port += "Type:any"
							} else {
								l.Port += "Type:" + strconv.FormatInt(*IpPermissionsEgress.FromPort, 10)
							}
							if *IpPermissionsEgress.ToPort == -1 {
								l.Port += "_Code:any"
							} else {
								l.Port += "_Code:" + strconv.FormatInt(*IpPermissionsEgress.ToPort, 10)
							}
						} else {
							l.Port = strconv.FormatInt(*IpPermissionsEgress.FromPort, 10)
						}

						if IpPermissionsEgress.IpRanges == nil {
							l.Address = "NULL"
						} else {
							for _, IpRanges := range IpPermissionsEgress.IpRanges {
								l.Address = *IpRanges.CidrIp
								ruleLists = append(ruleLists, l)
							}
						}

						if IpPermissionsEgress.UserIdGroupPairs == nil {
							l.Address = "NULL"
						} else {
							for _, UserIdGroupPairs := range IpPermissionsEgress.UserIdGroupPairs {
								l.Address = *UserIdGroupPairs.GroupId
								ruleLists = append(ruleLists, l)
							}
						}
					}
				}
			}
		}
	}
	return ruleLists
}

func injectRuleList(w *tabwriter.Writer, column string, ruleList []ruleList) {
	for _, l := range ruleList {
		fmt.Fprintf(
			w,
			column,
			l.Way,
			l.Protocol,
			l.Port,
			l.Address,
		)
	}
}
