package db

import (
	"fmt"
	"text/tabwriter"

	"ccmd/commons"
)

func printDbListColumn(w *tabwriter.Writer) {
	fmt.Fprintf(
		w,
		getDbListColumn(),
		"DB-ID",
		"STATE",
		"AZ (Primary)",
		"AZ (Secondary)",
		"MULTI-AZ",
		"PROVIDER",
	)
}

func printSnapListColumn(w *tabwriter.Writer) {
	fmt.Fprintf(
		w,
		getSnapListColumn(),
		"SNAPSHOT-ID",
		"DB-ID",
		"ENGINE",
		"AZ",
		"SIZE (GB)",
		"Progress (%)",
	)
}

func getDbListColumn() string {
	column := "%s"
	for i := 1; i < 6; i++ {
		column += commons.Flags.Delimiter + "%s"
	}
	column += "\n"
	return column
}

func getSnapListColumn() string {
	column := "%s"
	for i := 1; i < 4; i++ {
		column += commons.Flags.Delimiter + "%s"
	}
	column += "\n"
	return column
	// "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s" + listFlag.Delimiter + "%s\n"
}
