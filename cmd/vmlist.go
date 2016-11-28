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

var vmlistCmd = &cobra.Command{
	Use:   "list",
	Short: "VM list of all regions",
	Long:  "VM list of all regions",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup                                  //parallel processing counter group
		var instanceParamEC2 = []ec2.DescribeInstancesOutput{} //API response slide
		COLUMN := columnSettingVM()                            //define column
		regionsAWS := getAWSRegions()                          //get region list (AWS)

		w := new(tabwriter.Writer)          //generate tabwriter
		w.Init(os.Stdout, 0, 8, 2, '\t', 0) //configure tabwriter

		for _, region := range regionsAWS {
			wg.Add(1)                                      //waiting group count up
			go setEC2Param(region, &wg, &instanceParamEC2) //get instance information from API
			time.Sleep(1 * time.Millisecond)               //
		}
		wg.Wait()

		printColumn(COLUMN, w)
		printInstances(instanceParamEC2, COLUMN, w)

		w.Flush() //print

	},
}

func init() {
	vmCmd.AddCommand(vmlistCmd)
	vmlistCmd.Flags().StringVarP(&listFlag.Delimiter, "delimiter", "d", "\t", "delimiter") //define -d --delimiter flag
}

func printColumn(COLUMN string, w *tabwriter.Writer) {

	fmt.Fprintf(
		w,
		COLUMN,
		"NAME-TAG",
		"INSTANCE-ID",
		"STATE",
		"GLOBAL",
		"LOCAL",
		"AZ",
		"PROVIDER",
	)

}

func columnSettingVM() string {

	column := "%s"

	for i := 0; i < 6; i++ {
		column += listFlag.Delimiter + "%s"
	}
	column += "\n"

	return column

}

func printInstances(instanceParamEC2 []ec2.DescribeInstancesOutput, COLUMN string, w *tabwriter.Writer) {

	for _, lawJson := range instanceParamEC2 {
		for _, Reservations := range lawJson.Reservations {
			for _, Instances := range Reservations.Instances {

				var (
					Tags             string
					InstanceId       string
					State            string
					PublicIpAddress  string
					PrivateIpAddress string
					AvailabilityZone string
					Provider         string
				)

				if Instances.InstanceId == nil {
					InstanceId = "NULL"
				} else {
					InstanceId = *Instances.InstanceId
				}

				if Instances.State.Name == nil {
					State = "NULL"
				} else {
					State = *Instances.State.Name
				}

				if Instances.PublicIpAddress == nil {
					PublicIpAddress = "NULL"
				} else {
					PublicIpAddress = *Instances.PublicIpAddress
				}

				if Instances.PrivateIpAddress == nil {
					PrivateIpAddress = "NULL"
				} else {
					PrivateIpAddress = *Instances.PrivateIpAddress
				}

				if Instances.Placement.AvailabilityZone == nil {
					AvailabilityZone = "NULL"
				} else {
					AvailabilityZone = *Instances.Placement.AvailabilityZone
				}

				if Instances.Tags == nil {
					Tags = "NULL"
				} else if strings.Contains(*Instances.Tags[0].Value, listFlag.Delimiter) && (listFlag.Delimiter != "") { //When the delimiter is included in Name-Tag, enclose it with double quotation because separation will increase
					Tags = "\"" + *Instances.Tags[0].Value + "\"" //enclosed Name-Tag with double quotation
				} else {
					Tags = *Instances.Tags[0].Value
				}

				Provider = "AWS"

				fmt.Fprintf(
					w,
					COLUMN,
					Tags,
					InstanceId,
					State,
					PublicIpAddress,
					PrivateIpAddress,
					AvailabilityZone,
					Provider,
				)
			}
		}
	}

}
