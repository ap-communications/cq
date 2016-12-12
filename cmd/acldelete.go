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

var acldeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete ACL rule",
	Long: `delete ACL rule

Example:
  cq acl delete --groupid sg-fd8cc1ee --way ingress --protocol tcp --port 22 --address 192.0.2.0/24
`,
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup //parallel processing counter group
		var gid = []string{listFlag.GroupId}
		stats := map[string]int{}

		if checkSecurityGroupFlags() != "" {
			fmt.Println(checkSecurityGroupFlags())
			return
		}

		stats[listFlag.GroupId] = 0 //generate id hit counter (this function can't multi args)
		regionsAWS := getAWSRegions()

		for _, region := range regionsAWS {
			wg.Add(1) //waiting group count up
			go deleteSecurityGroupRule(region, &wg, stats)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		printHitId(gid, stats)

	},
}

func init() {
	aclCmd.AddCommand(acldeleteCmd)
	acldeleteCmd.Flags().StringVarP(&listFlag.GroupId, "groupid", "", "", "security group-id")                    // define --groupid flag
	acldeleteCmd.Flags().StringVarP(&listFlag.Protocol, "protocol", "", "", "tcp, udp, icmp, any (default: any)") // define --protocol flag
	acldeleteCmd.Flags().StringVarP(&listFlag.Address, "address", "", "", "CIDR address (default: 0.0.0.0/0)")    // define --address flag
	acldeleteCmd.Flags().StringVarP(&listFlag.Port, "port", "", "", "port (default: any)")                        // define --port flag
	acldeleteCmd.Flags().StringVarP(&listFlag.Way, "way", "", "", "ingress or egress")                            // define --way flag
}

func deleteSecurityGroupRule(region string, wg *sync.WaitGroup, stats map[string]int) {

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
				_, err := sginstance.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{ //add security group rule (ingress)
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
				_, err := sginstance.RevokeSecurityGroupEgress(&ec2.RevokeSecurityGroupEgressInput{ //add security group rule (egress)
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
