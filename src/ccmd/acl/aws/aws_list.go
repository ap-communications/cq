package aws

import (
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type aclList struct {
	GroupName   string
	Tags        string
	GroupId     string
	Description string
	Provider    string
}

func AclList(w *tabwriter.Writer, column string) {
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

	var sgs []ec2.DescribeSecurityGroupsOutput
	for _, region := range regions {
		sgs = append(sgs, m[region])
	}
	aclLists := getAclList(sgs)
	injectAclList(w, column, aclLists)
}

func getAclList(sgs []ec2.DescribeSecurityGroupsOutput) []aclList {
	var aclLists []aclList
	for _, lawJson := range sgs {
		for _, SecurityGroups := range lawJson.SecurityGroups {
			var l aclList

			if SecurityGroups.GroupName == nil {
				l.GroupName = "NULL"
			} else if strings.Contains(*SecurityGroups.GroupName, commons.Flags.Delimiter) && (commons.Flags.Delimiter != "") { //When the delimiter is included in GroupName, enclose it with double quotation because separation will increase
				l.GroupName = "\"" + *SecurityGroups.GroupName + "\"" //enclosed GroupName with double quotation
			} else {
				l.GroupName = *SecurityGroups.GroupName
			}

			if SecurityGroups.Tags == nil {
				l.Tags = "NULL"
			} else if strings.Contains(*SecurityGroups.Tags[0].Value, commons.Flags.Delimiter) && (commons.Flags.Delimiter != "") { //When the delimiter is included in Name-Tag, enclose it with double quotation because separation will increase
				l.Tags = "\"" + *SecurityGroups.Tags[0].Value + "\"" //enclosed Name-Tag with double quotation
			} else {
				l.Tags = *SecurityGroups.Tags[0].Value
			}

			if SecurityGroups.GroupId == nil {
				l.GroupId = "NULL"
			} else {
				l.GroupId = *SecurityGroups.GroupId
			}

			if SecurityGroups.Description == nil {
				l.Description = "NULL"
			} else if strings.Contains(*SecurityGroups.Description, commons.Flags.Delimiter) && (commons.Flags.Delimiter != "") { //When the delimiter is included in description, enclose it with double quotation because separation will increase
				l.Description = "\"" + *SecurityGroups.Description + "\"" //enclosed description with double quotation
			} else {
				l.Description = *SecurityGroups.Description
			}

			l.Provider = "AWS"

			aclLists = append(aclLists, l)
		}
	}
	return aclLists
}

func injectAclList(w *tabwriter.Writer, column string, aclLists []aclList) {
	for _, l := range aclLists {
		fmt.Fprintf(
			w,
			column,
			l.GroupName,
			l.Tags,
			l.GroupId,
			l.Description,
			l.Provider,
		)
	}
}
