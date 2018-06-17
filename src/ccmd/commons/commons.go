package commons

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	DEFAULT_REGION = "ap-northeast-1"
	VERSION        = "1.1.1"
)

var Flags struct {
	VersionFlag bool   // --version -v
	Delimiter   string // --delimiter -d
	ImageId     string // --imageid
	Region      string // --region
	Type        string // --type
	Keyname     string // --key
	Force       bool   // --force -f
	GroupId     string // --group-id
	Address     string // --address
	Protocol    string // --protocol
	Port        string // --port
	Way         string // --way
	NoFailover  bool   // --no-failover
	SnapshotId  string // --snapshot-id
	FilePath    string // --file
}

type InstanceList struct {
	Tags             string
	InstanceId       string
	State            string
	PublicIpAddress  string
	PrivateIpAddress string
	AvailabilityZone string
	Provider         string
}

// Get AWS regions list
func GetAwsRegions() []string {
	ec2Instance := ec2.New(session.New(), &aws.Config{Region: aws.String(DEFAULT_REGION)})
	resp, err := ec2Instance.DescribeRegions(nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	var awsRegions []string
	for _, r := range resp.Regions {
		awsRegions = append(awsRegions, *r.RegionName)
	}
	return awsRegions
}

// Confirmation before dangerous operation
func Confirm() bool {
	if Flags.Force { //if there is enabled force option, dont confirmation
		return true
	}
	var input string
	fmt.Printf("It may be impact of your service or data. Are you sure?  Y/N  ")
	fmt.Scanln(&input)
	if (input == "Y") || (input == "y") {
		return true
	}
	fmt.Printf("Cancelled\n")
	return false
}

// Check valid AWS region
func CheckAWSRegion(region string) bool {
	if region == "" { // if no set return true
		return true
	}
	regionList := GetAwsRegions()
	for _, r := range regionList {
		if r == region {
			return true
		}
	}
	return false
}
