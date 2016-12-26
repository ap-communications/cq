package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//flag variable
var listFlag struct {
	VersionFlag bool   // --version -v
	Delimiter   string // --delimiter -d
	GroupId     string // --groupid
	Address     string // --address
	Protocol    string // --protocol
	Port        string // --port
	Way         string // --way
	ImageId     string // --imageid
	Region      string // --region
	Type        string // --type
	Keyname     string // --key
	Force       bool   // --force -f
}

var dbFlag struct {
	Failover bool // --failover
}

//get region list (AWS)
func getAWSRegions() []string {

	var regionsAWS = []string{}

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")}) //generate API query instance
	resp, err := ec2instance.DescribeRegions(nil)                                            //get instance information from API

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, r1 := range resp.Regions {
		regionsAWS = append(regionsAWS, *r1.RegionName)
	}

	return regionsAWS
}

func printHitId(args []string, stats map[string]int) {

	for _, argid := range args { //if there is hit count 0, print "fail"
		if stats[argid] == 0 {
			fmt.Printf("Fail      %s   Not found\n", argid)
		} else if stats[argid] >= 2 { //if there is multi hit (count is over 1), show warning
			fmt.Printf("WARNING: Multi hit Instance-ID  %s", argid)
		}
	}

}

func checkRegion(region string) bool {

	regionList := getAWSRegions()

	//if there is valid AWS region or not set (use api-key default), return true
	for _, regionAWS := range regionList {
		if region == regionAWS {
			return true
		} else if region == "" {
			return true
		}
	}

	return false //if there is invalid region name, return false

}
