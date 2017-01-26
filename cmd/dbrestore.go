package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

var dbrestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore saved DB configure information file",
	Long:  "Restore saved DB configure information file",
	Run: func(cmd *cobra.Command, args []string) {

		if (dbFlag.SnapshotID == "") || (dbFlag.File == "") {
			fmt.Println("Missing required flag, must set --snapshot-id and --file")
			return
		}

		w := new(tabwriter.Writer)          //generate tabwriter
		w.Init(os.Stdout, 0, 8, 2, '\t', 0) //configure tabwriter

		restoreDBInstance(dbFlag.File, dbFlag.SnapshotID, w)

		w.Flush()

	},
}

func init() {
	dbCmd.AddCommand(dbrestoreCmd)
	dbrestoreCmd.Flags().StringVarP(&dbFlag.SnapshotID, "snapshot-id", "", "", "Snapshot ID")
	dbrestoreCmd.Flags().StringVarP(&dbFlag.File, "file", "", "", "Backuped configure file path")
}

func restoreDBInstance(filepath string, snapshotId string, w *tabwriter.Writer) {

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data rds.DBInstance
	unmarshalErr := json.Unmarshal([]byte(string(file)), &data)
	if unmarshalErr != nil {
		fmt.Println(err)
		return
	}

	region := *data.AvailabilityZone
	azSize := len(region) - 1 //AZ format is region name + one alphabet
	region = region[:azSize]

	inputDBName := ""
	fmt.Printf("Enter DB name\n")
	fmt.Scanln(&inputDBName) //stdin

	rdsInstance := rds.New(session.New(), &aws.Config{Region: aws.String(region)})
	_, errRestore := rdsInstance.RestoreDBInstanceFromDBSnapshot(&rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:         aws.String(*data.DBInstanceClass),
		DBInstanceIdentifier:    aws.String(inputDBName),
		DBSnapshotIdentifier:    aws.String(snapshotId),
		Engine:                  aws.String(*data.Engine),
		AutoMinorVersionUpgrade: aws.Bool(*data.AutoMinorVersionUpgrade),
		AvailabilityZone:        aws.String(*data.AvailabilityZone),
		CopyTagsToSnapshot:      aws.Bool(*data.CopyTagsToSnapshot),
		//		DBName:                  aws.String(*data.DBName),
		LicenseModel:       aws.String(*data.LicenseModel),
		MultiAZ:            aws.Bool(*data.MultiAZ),
		Port:               aws.Int64(*data.Endpoint.Port),
		PubliclyAccessible: aws.Bool(*data.PubliclyAccessible),
	})

	if errRestore != nil {
		fmt.Println(errRestore)
		return
	}

	fmt.Fprintf(w, "Started DB restore\n")

	return

}
