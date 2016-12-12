package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"
)

var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Show ACL rule",
	Long:  "Show ACL rule",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup            //parallel processing counter group
		stats := map[string]int{}        //instance id hit check map
		COLUMN := columnSettingACLRule() //define column
		regionsAWS := getAWSRegions()    //get region list (AWS)

		for _, argid := range args { //create hit judgment map for character string set as argument
			stats[argid] = 0 //init map (hit count is 0)
		}

		if len(args) == 0 { //If there is no argument, abort
			fmt.Printf("missing args (SecurityGroup-ID)\n")
			return //処理終了
		}

		w := new(tabwriter.Writer)          //generate tabwriter
		w.Init(os.Stdout, 0, 8, 2, '\t', 0) //configure tabwriter
		printRuleColumn(COLUMN, w)          //print column

		for _, region := range regionsAWS {
			wg.Add(1) //waiting group count up
			go printSecurityGroupRules(args, region, &wg, stats, w, COLUMN)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		w.Flush() //print
		printHitId(args, stats)

	},
}

func init() {
	aclCmd.AddCommand(ruleCmd)
}

func columnSettingACLRule() string {

	var column string

	if listFlag.Delimiter != "" {
		column = "%s"
		for i := 1; i < 4; i++ { //number of column is 4
			column += listFlag.Delimiter + "%s"
		}
		column += "\n"
	} else {
		column = "%s\t%s\t%s\t%s\t\n" // "%s\t%s\t"
	}

	return column

}

func printRuleColumn(COLUMN string, w *tabwriter.Writer) {

	fmt.Fprintf(
		w,
		COLUMN,
		"WAY",
		"PROTOCOL",
		"PORT",
		"ADDRESS",
	)

}

func printSecurityGroupRules(target []string, region string, wg *sync.WaitGroup, stats map[string]int, w *tabwriter.Writer, COLUMN string) {

	defer wg.Done()

	sgParamEC2 := getSecurityGroupParam(region)

	for _, SecurityGroups := range sgParamEC2.SecurityGroups {
		for _, sgid := range target {
			if *SecurityGroups.GroupId == sgid {

				stats[sgid]++

				for _, IpPermissions := range SecurityGroups.IpPermissions {

					var (
						Way      string //way of packet
						Protocol string
						Port     string
						Address  string
					)

					Way = "Ingress"

					if IpPermissions.IpProtocol == nil {
						Protocol = "NULL"
					} else if *IpPermissions.IpProtocol == "-1" {
						Protocol = "any"
					} else {
						Protocol = *IpPermissions.IpProtocol
					}

					if IpPermissions.FromPort == nil {
						Port = "any"
					} else if *IpPermissions.IpProtocol == "icmp" { //if protocol is icmp, format is Type & Code
						if *IpPermissions.FromPort == -1 {
							Port += "Type:any"
						} else {
							Port += "Type:" + strconv.FormatInt(*IpPermissions.FromPort, 10)
						}
						if *IpPermissions.ToPort == -1 {
							Port += "_Code:any"
						} else {
							Port += "_Code:" + strconv.FormatInt(*IpPermissions.ToPort, 10)
						}
					} else {
						Port = strconv.FormatInt(*IpPermissions.FromPort, 10)
					}

					if IpPermissions.IpRanges == nil {
						Address = "NULL"
					} else {
						for _, IpRanges := range IpPermissions.IpRanges {
							Address = *IpRanges.CidrIp
							fmt.Fprintf(w, COLUMN, Way, Protocol, Port, Address)
						}
					}

					if IpPermissions.UserIdGroupPairs == nil {
						Address = "NULL"
					} else {
						for _, UserIdGroupPairs := range IpPermissions.UserIdGroupPairs {
							Address = *UserIdGroupPairs.GroupId
							fmt.Fprintf(w, COLUMN, Way, Protocol, Port, Address)
						}
					}

				}

				for _, IpPermissionsEgress := range SecurityGroups.IpPermissionsEgress {

					var (
						Way      string //way of packet (ingress or egress)
						Protocol string
						Port     string
						Address  string
					)

					Way = "Egress"

					if IpPermissionsEgress.IpProtocol == nil {
						Protocol = "NULL"
					} else if *IpPermissionsEgress.IpProtocol == "-1" {
						Protocol = "any"
					} else {
						Protocol = *IpPermissionsEgress.IpProtocol
					}

					if IpPermissionsEgress.FromPort == nil {
						Port = "any"
					} else if *IpPermissionsEgress.IpProtocol == "icmp" { //if protocol is icmp, format is Type & Code
						if *IpPermissionsEgress.FromPort == -1 {
							Port += "Type:any"
						} else {
							Port += "Type:" + strconv.FormatInt(*IpPermissionsEgress.FromPort, 10)
						}
						if *IpPermissionsEgress.ToPort == -1 {
							Port += "_Code:any"
						} else {
							Port += "_Code:" + strconv.FormatInt(*IpPermissionsEgress.ToPort, 10)
						}
					} else {
						Port = strconv.FormatInt(*IpPermissionsEgress.FromPort, 10)
					}

					if IpPermissionsEgress.IpRanges == nil {
						Address = "NULL"
					} else {
						for _, IpRanges := range IpPermissionsEgress.IpRanges {
							Address = *IpRanges.CidrIp
							fmt.Fprintf(w, COLUMN, Way, Protocol, Port, Address)
						}
					}

					if IpPermissionsEgress.UserIdGroupPairs == nil {
						Address = "NULL"
					} else {
						for _, UserIdGroupPairs := range IpPermissionsEgress.UserIdGroupPairs {
							Address = *UserIdGroupPairs.GroupId
							fmt.Fprintf(w, COLUMN, Way, Protocol, Port, Address)
						}
					}

				}

			}
		}
	}

	return

}
