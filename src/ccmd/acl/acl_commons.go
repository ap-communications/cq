package acl

import (
	"fmt"
	"text/tabwriter"

	"github.com/ap-communications/cq/src/ccmd/commons"
)

func printAclListColumn(w *tabwriter.Writer) {
	fmt.Fprintf(
		w,
		getAclListColumn(),
		"GROUP-NAME",
		"NAME-TAG",
		"ID",
		"DESCRIPTION",
		"PROVIDER",
	)
}

func printRuleListColumn(w *tabwriter.Writer) {
	fmt.Fprintf(
		w,
		getRuleListColumn(),
		"WAY",
		"PROTOCOL",
		"PORT",
		"ADDRESS",
	)
}

func getAclListColumn() string {
	var column string
	if commons.Flags.Delimiter != "" {
		column = "%s"
		for i := 1; i < 5; i++ {
			column += commons.Flags.Delimiter + "%s"
		}
		column += "\n"
	} else {
		column = "%s\t%s\t%s\t%s\t%s\t\n"
	}
	return column
}

func getRuleListColumn() string {
	var column string
	if commons.Flags.Delimiter != "" {
		column = "%s"
		for i := 1; i < 4; i++ {
			column += commons.Flags.Delimiter + "%s"
		}
		column += "\n"
	} else {
		column = "%s\t%s\t%s\t%s\t\n"
	}
	return column
}
