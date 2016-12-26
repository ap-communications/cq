package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

var dbsnaplistCmd = &cobra.Command{
	Use:   "snaplist",
	Short: "DB snapshot list of all regions",
	Long:  "DB snapshot list of all regions",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup //parallel processing counter group
		var snapshotParamRDS = []rds.DescribeDBSnapshotsOutput{}
		COLUMN := columnSettingDBSnapshot()
		regionsAWS := getAWSRegions()

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)

		for _, region := range regionsAWS {
			wg.Add(1)
			go setRDSSnapList(region, &wg, &snapshotParamRDS)
			time.Sleep(1 * time.Millisecond) //switch goroutine
		}
		wg.Wait()

		printDBSnapshotColumn(w)
		printDBSnapshots(snapshotParamRDS, COLUMN, w)

		w.Flush()

	},
}

func init() {
	dbCmd.AddCommand(dbsnaplistCmd)
	dbsnaplistCmd.Flags().StringVarP(&listFlag.Delimiter, "delimiter", "d", "\t", "delimiter") //define -d --delimiter flag
}

func printDBSnapshotColumn(w *tabwriter.Writer) {

	column := "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s\n"

	fmt.Fprintf(
		w,
		column,
		"SNAPSHOT-ID",
		"DB-ID",
		"ENGINE",
		"AZ",
		"SIZE (GB)",
		"Progress (%)",
	)

}

func columnSettingDBSnapshot() string {

	column := "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%d" + listFlag.Delimiter + "%d\n"

	return column

}

func printDBSnapshots(snapshotParamRDS []rds.DescribeDBSnapshotsOutput, COLUMN string, w *tabwriter.Writer) {

	for _, lawJson := range snapshotParamRDS {
		for _, DBSnapshots := range lawJson.DBSnapshots {

			var (
				SnapshotID       string
				DBID             string
				Engine           string
				AvailabilityZone string
				Size             int64
				Progress         int64
			)

			if DBSnapshots.DBSnapshotIdentifier == nil {
				SnapshotID = "NULL"
			} else {
				SnapshotID = *DBSnapshots.DBSnapshotIdentifier
			}

			if DBSnapshots.DBInstanceIdentifier == nil {
				DBID = "NULL"
			} else {
				DBID = *DBSnapshots.DBInstanceIdentifier
			}

			if DBSnapshots.Engine == nil {
				Engine = "NULL"
			} else {
				Engine = *DBSnapshots.Engine + "-" + *DBSnapshots.EngineVersion
			}

			if DBSnapshots.AvailabilityZone == nil {
				AvailabilityZone = "NULL"
			} else {
				AvailabilityZone = *DBSnapshots.AvailabilityZone
			}

			if DBSnapshots.AllocatedStorage == nil {
				Size = 0
			} else {
				Size = *DBSnapshots.AllocatedStorage
			}

			if DBSnapshots.PercentProgress == nil {
				Progress = 0
			} else {
				Progress = *DBSnapshots.PercentProgress
			}

			fmt.Fprintf(
				w,
				COLUMN,
				SnapshotID,
				DBID,
				Engine,
				AvailabilityZone,
				Size,
				Progress,
			)
		}
	}

}
