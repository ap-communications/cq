package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
)

var vmeasyupCmd = &cobra.Command{
	Use:   "easyup",
	Short: "Create and run new VM",
	Long:  "Create and run new VM",
	Run: func(cmd *cobra.Command, args []string) {

		if checkRegion(listFlag.Region) != true {
			fmt.Printf("Invalid region\n[region list]\n")
			fmt.Printf("%s", getAWSRegions)
			return
		}

		createInstance()

	},
}

func init() {
	vmCmd.AddCommand(vmeasyupCmd)
	vmeasyupCmd.Flags().StringVarP(&listFlag.Region, "region", "", "ap-northeast-1", "Region")
	vmeasyupCmd.Flags().StringVarP(&listFlag.ImageId, "imageid", "", "", "Instance image ID (default latest Amazon Linux image)")
	vmeasyupCmd.Flags().StringVarP(&listFlag.Type, "type", "", "t2.micro", "Instance type")
	vmeasyupCmd.Flags().StringVarP(&listFlag.Keyname, "key", "", "", "SSH key-pair (default new generate)")
	vmeasyupCmd.Flags().StringVarP(&listFlag.GroupId, "groupid", "", "", "Security Group ID (default new generate)")
}

func createInstance() {

	region := listFlag.Region
	imageId := listFlag.ImageId
	iType := listFlag.Type
	keyname := listFlag.Keyname
	sgid := listFlag.GroupId
	qt, _ := strconv.ParseInt("1", 10, 64) //static 1 vm

	var key string
	if keyname == "" {
		key, keyname = createKey(region) //default key generate
	}

	if sgid == "" {
		sgid = createSecurityGroup(region) //default security group generate
	}

	if imageId == "" { //default image-id is latest Amazon Linux
		imageId = getLatestAmazonLinuxInstance(region)
	}

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := ec2instance.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(imageId),
		MaxCount:     aws.Int64(qt),
		MinCount:     aws.Int64(qt),
		InstanceType: aws.String(iType),
		KeyName:      aws.String(keyname),
		SecurityGroupIds: []*string{
			aws.String(sgid),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	checkedResp := checkInstanceCreated(region, *resp.Instances[0].InstanceId)

	_, tagErr := ec2instance.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{ aws.String(*resp.Instances[0].InstanceId) },
		Tags: []*ec2.Tag{
			{
				Key: aws.String("Name"),
				Value: aws.String("cq-easyup"),
			},
		},
	})
	if tagErr != nil {
		fmt.Println(tagErr)
	}

	fmt.Printf("     Instance ID: %s\n", *checkedResp.Reservations[0].Instances[0].InstanceId)

	if listFlag.GroupId == "" { //new generated (not designate security group id) ssh port is permit from 0.0.0.0/0
		addSSHFree(region, sgid)
		fmt.Printf("SecurityGroup ID: %s\n", sgid)
		fmt.Printf("\n  ***** IMPORTANT: SSH (TCP22) is anyone can access!! *****\n\n")
	} else {
		fmt.Printf("SecurityGroup ID: %s\n", sgid)
	}

	fmt.Printf("          Global: %s\n", *checkedResp.Reservations[0].Instances[0].PublicIpAddress)

	if listFlag.Keyname == "" { //there is not designate ssh keypair, generate new keypair
		fmt.Printf("         SSH Key:\n%s\n", key)
	} else {
		fmt.Printf("         SSH Key: %s\n", keyname)
	}

	return

}

func getLatestAmazonLinuxInstance(region string) string {

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //generate API query instance
	resp, err := ec2instance.DescribeImages(&ec2.DescribeImagesInput{
		Owners: []*string{
			aws.String("amazon"),
		},
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("architecture"),
				Values: []*string{aws.String("x86_64")},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []*string{aws.String("hvm")},
			},
			{
				Name:   aws.String("root-device-type"),
				Values: []*string{aws.String("ebs")},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
			{
				Name:   aws.String("description"),
				Values: []*string{aws.String("Amazon Linux AMI*")},
			},
			{
				Name:   aws.String("block-device-mapping.volume-type"),
				Values: []*string{aws.String("gp2")},
			},
			{
				Name:   aws.String("image-type"),
				Values: []*string{aws.String("machine")},
			},
			{
				Name:   aws.String("block-device-mapping.volume-size"),
				Values: []*string{aws.String("8")},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	//exclusion ecs-instance, nat-instance, rc-version
	ids := map[string]string{}
	for _, Images := range resp.Images {
		if (strings.Contains(*Images.Name, "ecs") || strings.Contains(*Images.Name, "nat") || strings.Contains(*Images.Name, "rc")) != true {
			ids[*Images.CreationDate] = *Images.ImageId
		}
	}

	//sort by date
	var keys []string
	for k, _ := range ids {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return ids[keys[len(keys)-1]] //latest date image-id

}

func addSSHFree(region string, sgid string) {

	var port int64 = 22    //SSH
	protocol := "TCP"      //SSH
	address := "0.0.0.0/0" //free

	sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})               //create ec2(security group) api-instance
	_, err := sginstance.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{ //add security group rule (ingress)
		GroupId: aws.String(sgid),
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
	})

	if err != nil { //if got error, print it
		fmt.Println(err)
		return
	}

}
