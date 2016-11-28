package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

var rulelistCmd = &cobra.Command{
	Use:   "list",
	Short: "List acl",
	Long:  "List acl",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup                                 //parallel processing counter group
		var sgParamEC2 = []ec2.DescribeSecurityGroupsOutput{} //API response slide
		COLUMN := columnSettingACLList()                      //define column
		regionsAWS := getAWSRegions()                         //get region list (AWS)

		w := new(tabwriter.Writer)          //generate tabwriter
		w.Init(os.Stdout, 0, 8, 2, '\t', 0) //configure tabwriter

		for _, region := range regionsAWS {
			wg.Add(1)                                          //waiting group count up
			go setSecurityGroupParam(region, &wg, &sgParamEC2) //get instance information from API
			time.Sleep(1 * time.Millisecond)                   //
		}
		wg.Wait()

		printACLColumn(COLUMN, w)
		printSecurityGroups(sgParamEC2, COLUMN, w)

		w.Flush() //print

	},
}

func init() {
	aclCmd.AddCommand(rulelistCmd)
	rulelistCmd.Flags().StringVarP(&listFlag.Delimiter, "delimiter", "d", "", "delimiter(default:tab)") //define -d --delimiter flag
}

func columnSettingACLList() string {

	var column string

	if listFlag.Delimiter != "" {
		column = "%s"
		for i := 1; i < 5; i++ { //number of column is 5
			column += listFlag.Delimiter + "%s"
		}
		column += "\n"
	} else {
		column = "%s\t%s\t%s\t%s\t%s\t\n"
	}

	return column

}

func printACLColumn(COLUMN string, w *tabwriter.Writer) {

	fmt.Fprintf(
		w,
		COLUMN,
		"GROUP-NAME",
		"NAME-TAG",
		"ID",
		"DESCRIPTION",
		"PROVIDER",
	)

}

func printSecurityGroups(sgParamEC2 []ec2.DescribeSecurityGroupsOutput, COLUMN string, w *tabwriter.Writer) {

	for _, lawJson := range sgParamEC2 {
		for _, SecurityGroups := range lawJson.SecurityGroups {

			var (
				GroupName   string
				Tags        string
				GroupId     string
				Description string
				Provider    string
			)

			if SecurityGroups.GroupName == nil {
				GroupName = "NULL"
			} else if strings.Contains(*SecurityGroups.GroupName, listFlag.Delimiter) && (listFlag.Delimiter != "") { //When the delimiter is included in GroupName, enclose it with double quotation because separation will increase
				GroupName = "\"" + *SecurityGroups.GroupName + "\"" //enclosed GroupName with double quotation
			} else {
				GroupName = *SecurityGroups.GroupName
			}

			if SecurityGroups.Tags == nil {
				Tags = "NULL"
			} else if strings.Contains(*SecurityGroups.Tags[0].Value, listFlag.Delimiter) && (listFlag.Delimiter != "") { //When the delimiter is included in Name-Tag, enclose it with double quotation because separation will increase
				Tags = "\"" + *SecurityGroups.Tags[0].Value + "\"" //enclosed Name-Tag with double quotation
			} else {
				Tags = *SecurityGroups.Tags[0].Value
			}

			if SecurityGroups.GroupId == nil {
				GroupId = "NULL"
			} else {
				GroupId = *SecurityGroups.GroupId
			}

			if SecurityGroups.Description == nil {
				Description = "NULL"
			} else if strings.Contains(*SecurityGroups.Description, listFlag.Delimiter) && (listFlag.Delimiter != "") { //When the delimiter is included in description, enclose it with double quotation because separation will increase
				Description = "\"" + *SecurityGroups.Description + "\"" //enclosed description with double quotation
			} else {
				Description = *SecurityGroups.Description
			}

			Provider = "AWS"

			fmt.Fprintf(
				w,
				COLUMN,
				GroupName,
				Tags,
				GroupId,
				Description,
				Provider,
			)

		}
	}

}
