package db

import (
	"os"
	"text/tabwriter"

	"github.com/ap-communications/cq/src/ccmd/db/aws"
)

var w = tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

func Destroy(args []string) {
	aws.Destroy(args)
}

func Inspect(args []string) {
	aws.Inspect(args)
}

func List() {
	printDbListColumn(w)
	aws.List(w, getDbListColumn())
	w.Flush()
}

func Reboot(args []string) {
	aws.Reboot(args)
}

func SnapDel(args []string) {
	aws.SnapDel(args)
}

func SnapList() {
	printSnapListColumn(w)
	aws.SnapList(w, getSnapListColumn())
	w.Flush()
}

func Snapshot(args []string) {
	aws.Snapshot(args)
}
