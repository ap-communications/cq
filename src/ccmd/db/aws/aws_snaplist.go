package aws

import (
	"fmt"
	"sync"
	"text/tabwriter"

	"ccmd/commons"

	"github.com/aws/aws-sdk-go/service/rds"
)

type snapList struct {
	SnapshotID       string
	DBID             string
	Engine           string
	AvailabilityZone string
	Size             int64
	Progress         int64
}

func SnapList(w *tabwriter.Writer, column string) {
	regions := commons.GetAwsRegions()
	m := map[string]rds.DescribeDBSnapshotsOutput{}
	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			i, err := getRdsSnapList(region)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			m[region] = i
		}(region)
	}
	wg.Wait()

	var rdsInstances []rds.DescribeDBSnapshotsOutput
	for _, region := range regions {
		rdsInstances = append(rdsInstances, m[region])
	}
	injectSnapList(w, column, getSnapList(rdsInstances))
}

func getSnapList(snapshotParamRDS []rds.DescribeDBSnapshotsOutput) []snapList {
	var snapLists []snapList
	for _, lawJson := range snapshotParamRDS {
		for _, DBSnapshots := range lawJson.DBSnapshots {
			var l snapList

			if DBSnapshots.DBSnapshotIdentifier == nil {
				l.SnapshotID = "NULL"
			} else {
				l.SnapshotID = *DBSnapshots.DBSnapshotIdentifier
			}

			if DBSnapshots.DBInstanceIdentifier == nil {
				l.DBID = "NULL"
			} else {
				l.DBID = *DBSnapshots.DBInstanceIdentifier
			}

			if DBSnapshots.Engine == nil {
				l.Engine = "NULL"
			} else {
				l.Engine = *DBSnapshots.Engine + "-" + *DBSnapshots.EngineVersion
			}

			if DBSnapshots.AvailabilityZone == nil {
				l.AvailabilityZone = "NULL"
			} else {
				l.AvailabilityZone = *DBSnapshots.AvailabilityZone
			}

			if DBSnapshots.AllocatedStorage == nil {
				l.Size = 0
			} else {
				l.Size = *DBSnapshots.AllocatedStorage
			}

			if DBSnapshots.PercentProgress == nil {
				l.Progress = 0
			} else {
				l.Progress = *DBSnapshots.PercentProgress
			}

			snapLists = append(snapLists, l)
		}
	}
	return snapLists
}

func injectSnapList(w *tabwriter.Writer, column string, snapLists []snapList) {
	for _, l := range snapLists {
		fmt.Fprintf(
			w,
			column,
			l.SnapshotID,
			l.DBID,
			l.Engine,
			l.AvailabilityZone,
			l.Size,
			l.Progress,
		)
	}
}
