package vm

import (
	"fmt"
	"os"
	"text/tabwriter"

	"ccmd/commons"
	"ccmd/vm/aws"
)

func List() {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	printInstanceListColumn(w)

	awsCh := make(chan commons.InstanceList)
	go aws.List(awsCh)

	var awsList []commons.InstanceList
	for {
		select {
		case awsInstance, ok := <-awsCh:
			if ok {
				awsList = append(awsList, awsInstance)
			} else {
				awsCh = nil
				inject(w, awsList)
			}
		}
		if awsCh == nil { // If there is all finished, print and return.
			w.Flush()
			return
		}
	}
}

func Start(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.StartInstance(args)
}

func Stop(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.StopInstance(args)
}

func Reboot(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.RebootInstance(args)
}

func Destroy(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.DestroyInstance(args)
}

func Backup(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.BackupInstance(args)
}

func Restore(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.RestoreInstance(args)
}

func Easyup() {
	if commons.CheckAWSRegion(commons.Flags.Region) {
		aws.Easyup()
		return
	}
	fmt.Printf("Invalid region\n%v\n", commons.GetAwsRegions())
}

func Inspect(args []string) {
	if len(args) == 0 {
		fmt.Printf("missing args (Instance-ID)\n")
		return
	}
	aws.Inspect(args)
}
