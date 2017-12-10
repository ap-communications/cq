package vm

import (
	"fmt"
	"text/tabwriter"

	"ccmd/commons"
)

func inject(w *tabwriter.Writer, instanceLists []commons.InstanceList) {
	for _, l := range instanceLists {
		fmt.Fprintf(
			w,
			getColumn(),
			l.Tags,
			l.InstanceId,
			l.State,
			l.PublicIpAddress,
			l.PrivateIpAddress,
			l.AvailabilityZone,
			l.Provider,
		)
	}
}

func printInstanceListColumn(w *tabwriter.Writer) {
	fmt.Fprintf(
		w,
		getColumn(),
		"NAME-TAG",
		"INSTANCE-ID",
		"STATE",
		"GLOBAL",
		"LOCAL",
		"AZ",
		"PROVIDER",
	)
}

func getColumn() string {
	column := "%s"
	for i := 0; i < 6; i++ {
		column += commons.Flags.Delimiter + "%s"
	}
	return column + "\n"
}
