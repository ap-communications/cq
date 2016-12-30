package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

var dblistCmd = &cobra.Command{
	Use:   "list",
	Short: "DB list of all regions",
	Long:  "DB list of all regions",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup //parallel processing counter group
		var instanceParamRDS = []rds.DescribeDBInstancesOutput{}
		COLUMN := columnSettingDB()
		regionsAWS := getAWSRegions()

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)

		for _, region := range regionsAWS {
			wg.Add(1)
			go setRDSParam(region, &wg, &instanceParamRDS)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printDBColumn(COLUMN, w)
		printDBs(instanceParamRDS, COLUMN, w)

		w.Flush()

	},
}

func init() {
	dbCmd.AddCommand(dblistCmd)
	dblistCmd.Flags().StringVarP(&listFlag.Delimiter, "delimiter", "d", "\t", "delimiter") //define -d --delimiter flag
}

func printDBColumn(COLUMN string, w *tabwriter.Writer) {

	fmt.Fprintf(
		w,
		COLUMN,
		"DB-ID",
		"STATE",
		"AZ (Primary)",
		"AZ (Secondary)",
		"MULTI-AZ",
		"PROVIDER",
	)

}

func columnSettingDB() string {

	column := "%s"

	for i := 1; i < 6; i++ {
		column += listFlag.Delimiter + "%s"
	}
	column += "\n"

	return column

}

func printDBs(instanceParamRDS []rds.DescribeDBInstancesOutput, COLUMN string, w *tabwriter.Writer) {

	for _, lawJson := range instanceParamRDS {
		for _, DBInstances := range lawJson.DBInstances {

			var (
				DBId                      string
				State                     string
				AvailabilityZonePrimary   string
				AvailabilityZoneSecondary string
				MultiAz                   string
				Provider                  string
			)

			if DBInstances.DBInstanceIdentifier == nil {
				DBId = "NULL"
			} else if strings.Contains(*DBInstances.DBInstanceIdentifier, listFlag.Delimiter) && (listFlag.Delimiter != "") { //When the delimiter is included in Name-Tag, enclose it with double quotation because separation will increase
				DBId = "\"" + *DBInstances.DBInstanceIdentifier + "\"" //enclosed DBID with double quotation
			} else {
				DBId = *DBInstances.DBInstanceIdentifier
			}

			if DBInstances.DBInstanceStatus == nil {
				State = "NULL"
			} else {
				State = *DBInstances.DBInstanceStatus
			}

			if DBInstances.AvailabilityZone == nil {
				AvailabilityZonePrimary = "NULL"
			} else {
				AvailabilityZonePrimary = *DBInstances.AvailabilityZone
			}

			if DBInstances.SecondaryAvailabilityZone == nil {
				AvailabilityZoneSecondary = "NULL"
			} else {
				AvailabilityZoneSecondary = *DBInstances.SecondaryAvailabilityZone
			}

			if DBInstances.MultiAZ == nil {
				MultiAz = "NULL"
			} else if *DBInstances.MultiAZ == true {
				MultiAz = "Yes"
			} else if *DBInstances.MultiAZ == false {
				MultiAz = "No"
			} else {
				MultiAz = "Unknown"
			}

			Provider = "AWS"

			fmt.Fprintf(
				w,
				COLUMN,
				DBId,
				State,
				AvailabilityZonePrimary,
				AvailabilityZoneSecondary,
				MultiAz,
				Provider,
			)
		}
	}

}
