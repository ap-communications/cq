package aws

import (
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/aws/aws-sdk-go/service/rds"
)

type dbList struct {
	DbId                      string
	State                     string
	AvailabilityZonePrimary   string
	AvailabilityZoneSecondary string
	MultiAz                   string
	Provider                  string
}

func List(w *tabwriter.Writer, column string) {
	m := map[string]rds.DescribeDBInstancesOutput{}
	regions := commons.GetAwsRegions()
	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			i, err := getRdsInstances(region)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			m[region] = i
		}(region)
	}
	wg.Wait()

	var rdsInstances []rds.DescribeDBInstancesOutput
	for _, region := range regions {
		rdsInstances = append(rdsInstances, m[region])
	}
	injectDbList(w, column, getDbList(rdsInstances))
}

func getDbList(rdsInstances []rds.DescribeDBInstancesOutput) []dbList {
	var dbLists []dbList
	for _, rdsInstance := range rdsInstances {
		for _, DBInstances := range rdsInstance.DBInstances {
			var l dbList

			if DBInstances.DBInstanceIdentifier == nil {
				l.DbId = "NULL"
			} else if strings.Contains(*DBInstances.DBInstanceIdentifier, commons.Flags.Delimiter) && (commons.Flags.Delimiter != "") { //When the delimiter is included in Name-Tag, enclose it with double quotation because separation will increase
				l.DbId = "\"" + *DBInstances.DBInstanceIdentifier + "\"" //enclosed DbId with double quotation
			} else {
				l.DbId = *DBInstances.DBInstanceIdentifier
			}

			if DBInstances.DBInstanceStatus == nil {
				l.State = "NULL"
			} else {
				l.State = *DBInstances.DBInstanceStatus
			}

			if DBInstances.AvailabilityZone == nil {
				l.AvailabilityZonePrimary = "NULL"
			} else {
				l.AvailabilityZonePrimary = *DBInstances.AvailabilityZone
			}

			if DBInstances.SecondaryAvailabilityZone == nil {
				l.AvailabilityZoneSecondary = "NULL"
			} else {
				l.AvailabilityZoneSecondary = *DBInstances.SecondaryAvailabilityZone
			}

			if DBInstances.MultiAZ == nil {
				l.MultiAz = "NULL"
			} else if *DBInstances.MultiAZ == true {
				l.MultiAz = "Yes"
			} else if *DBInstances.MultiAZ == false {
				l.MultiAz = "No"
			} else {
				l.MultiAz = "Unknown"
			}

			l.Provider = "AWS"

			dbLists = append(dbLists, l)
		}
	}
	return dbLists
}

func injectDbList(w *tabwriter.Writer, column string, dbLists []dbList) {
	for _, l := range dbLists {
		fmt.Fprintf(
			w,
			column,
			l.DbId,
			l.State,
			l.AvailabilityZonePrimary,
			l.AvailabilityZoneSecondary,
			l.MultiAz,
			l.Provider,
		)
	}
}
