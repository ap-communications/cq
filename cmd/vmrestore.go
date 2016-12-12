package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

var vmrestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore saved configure information file",
	Long:  "Restore saved configure information file",
	Run: func(cmd *cobra.Command, args []string) {

		var wg sync.WaitGroup //parallel processing counter group

		if len(args) == 0 { //If there is no argument, abort
			fmt.Printf("missing args (file)\n")
			return
		}

		w := new(tabwriter.Writer)          //generate tabwriter
		w.Init(os.Stdout, 0, 8, 2, '\t', 0) //configure tabwriter

		for _, filepath := range args {
			wg.Add(1) //waiting group count up
			go restoreInstance(filepath, &wg, w)
			time.Sleep(1 * time.Millisecond) //
		}
		wg.Wait()

		w.Flush()

	},
}

func init() {

	USAGE := `Usage:
  cq vm restore [filename] [filename] ...
`

	vmCmd.AddCommand(vmrestoreCmd)
	vmrestoreCmd.SetUsageTemplate(USAGE) //override Usage words
}

func restoreInstance(filepath string, wg *sync.WaitGroup, w *tabwriter.Writer) {

	defer wg.Done()

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data ec2.Instance
	unmarshalErr := json.Unmarshal([]byte(string(file)), &data)
	if unmarshalErr != nil {
		fmt.Println(err)
		return
	}

	region := *data.Placement.AvailabilityZone
	azSize := len(region) - 1 //AZ format is region name + one alphabet
	region = region[:azSize]

	var monitoringEnabled bool
	if *data.Monitoring.State == "disabled" {
		monitoringEnabled = false
	} else {
		monitoringEnabled = true
	}

	ec2instance := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	resp, err := ec2instance.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(*data.ImageId),
		MaxCount:     aws.Int64(1), //static
		MinCount:     aws.Int64(1), //static
		EbsOptimized: aws.Bool(*data.EbsOptimized),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(*data.IamInstanceProfile.Arn),
		},
		InstanceType: aws.String(*data.InstanceType),
		KeyName:      aws.String(*data.KeyName),
		Monitoring: &ec2.RunInstancesMonitoringEnabled{
			Enabled: aws.Bool(monitoringEnabled),
		},
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(*data.Placement.AvailabilityZone),
			Tenancy:          aws.String(*data.Placement.Tenancy),
		},
		SecurityGroupIds: []*string{
			aws.String(*data.SecurityGroups[0].GroupId),
		},
		SubnetId: aws.String(*data.SubnetId),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	checkedResp := checkInstanceCreated(region, *resp.Instances[0].InstanceId)

	fmt.Fprintf(w, "     Instance ID: %s\n", *checkedResp.Reservations[0].Instances[0].InstanceId)
	fmt.Fprintf(w, "          Global: %s\n\n", *checkedResp.Reservations[0].Instances[0].PublicIpAddress)

	return

}
