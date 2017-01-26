package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var acladdCmd = &cobra.Command{
	Use:   "add",
	Short: "add ACL rule",
	Long: `add ACL rule

Example:
  cq acl add --groupid sg-fd8cc1ee --way ingress --protocol tcp --port 22 --address 192.0.2.0/24
`,
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup //parallel processing counter group
		var gid = []string{listFlag.GroupId}
		stats := map[string]int{} //define id hit counter map

		if checkSecurityGroupFlags() != "" {
			fmt.Println(checkSecurityGroupFlags())
			return
		}

		stats[listFlag.GroupId] = 0 //generate id hit counter (this function can't multi args)
		regionsAWS := getAWSRegions()

		for _, region := range regionsAWS {
			wg.Add(1) //waiting group count up
			go addSecurityGroupRule(region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		printHitId(gid, stats)

	},
}

func init() {
	aclCmd.AddCommand(acladdCmd)
	acladdCmd.Flags().StringVarP(&listFlag.GroupId, "groupid", "", "", "security group-id")                    // define --groupid flag
	acladdCmd.Flags().StringVarP(&listFlag.Protocol, "protocol", "", "", "tcp, udp, icmp, any (default: any)") // define --protocol flag
	acladdCmd.Flags().StringVarP(&listFlag.Address, "address", "", "", "CIDR address (default: 0.0.0.0/0)")    // define --address flag
	acladdCmd.Flags().StringVarP(&listFlag.Port, "port", "", "", "port (default: any)")                        // define --port flag
	acladdCmd.Flags().StringVarP(&listFlag.Way, "way", "", "", "ingress or egress")                            // define --way flag
}

func addSecurityGroupRule(region string, wg *sync.WaitGroup, stats map[string]int) {

	defer wg.Done()

	sgParamEC2 := getSecurityGroupParam(region)                   //get security group parameter
	port := translateSecurityGroupPort(listFlag.Port)             //translate input port number
	protocol := translateSecurityGroupProtocol(listFlag.Protocol) //translate input protocol string
	address := translateSecurityGroupAddress(listFlag.Address)    //translate input address string

	for _, SecurityGroups := range sgParamEC2.SecurityGroups {
		if *SecurityGroups.GroupId == listFlag.GroupId {
			stats[listFlag.GroupId]++                                                     //increment hit id counter
			sginstance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)}) //create ec2(security group) api-instance
			if listFlag.Way == "ingress" {
				_, err := sginstance.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{ //add security group rule (ingress)
					GroupId: aws.String(listFlag.GroupId),
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
			} else if listFlag.Way == "egress" {
				_, err := sginstance.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{ //add security group rule (egress)
					GroupId: aws.String(listFlag.GroupId),
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
			fmt.Printf("Success!\n")
		}
	}

	return

}
