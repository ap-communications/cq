package vm

import (
	"fmt"
	"text/tabwriter"

	"ccmd/commons"
)

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
