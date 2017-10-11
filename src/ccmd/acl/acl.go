package acl

import (
	"os"
	"text/tabwriter"

	"ccmd/acl/aws"
)

var w = tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

func Add(args []string) {
	aws.Add(args)
}

func Delete(args []string) {
	aws.Delete(args)
}

func Destroy(args []string) {
	aws.Destroy(args)
}

func AclList() {
	printAclListColumn(w)
	aws.AclList(w, getAclListColumn())
	w.Flush()
}

func RuleList(args []string) {
	printRuleListColumn(w)
	aws.RuleList(w, getRuleListColumn(), args)
	w.Flush()
}
