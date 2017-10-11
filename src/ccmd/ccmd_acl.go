package ccmd

import (
	"ccmd/acl"
	"ccmd/commons"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(aclCmd)

	aclCmd.AddCommand(aclAddCmd)
	aclAddCmd.Flags().StringVarP(&commons.Flags.GroupId, "group-id", "", "", "security group-id")
	aclAddCmd.Flags().StringVarP(&commons.Flags.Protocol, "protocol", "", "", "tcp, udp, icmp, any (default: any)")
	aclAddCmd.Flags().StringVarP(&commons.Flags.Address, "address", "", "", "CIDR address (default: 0.0.0.0/0)")
	aclAddCmd.Flags().StringVarP(&commons.Flags.Port, "port", "", "", "port (default: any)")
	aclAddCmd.Flags().StringVarP(&commons.Flags.Way, "way", "", "", "ingress or egress")

	aclCmd.AddCommand(aclDeleteCmd)
	aclDeleteCmd.Flags().StringVarP(&commons.Flags.GroupId, "group-id", "", "", "security group-id")
	aclDeleteCmd.Flags().StringVarP(&commons.Flags.Protocol, "protocol", "", "", "tcp, udp, icmp, any (default: any)")
	aclDeleteCmd.Flags().StringVarP(&commons.Flags.Address, "address", "", "", "CIDR address (default: 0.0.0.0/0)")
	aclDeleteCmd.Flags().StringVarP(&commons.Flags.Port, "port", "", "", "port (default: any)")
	aclDeleteCmd.Flags().StringVarP(&commons.Flags.Way, "way", "", "", "ingress or egress")

	aclCmd.AddCommand(aclDestroyCmd)
	aclDestroyCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "Destroy without confirmation")

	aclCmd.AddCommand(ruleListCmd)
	ruleListCmd.Flags().StringVarP(&commons.Flags.Delimiter, "delimiter", "d", "", "delimiter(default:tab)")

	aclCmd.AddCommand(ruleCmd)
}

var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "Control network access",
	Long:  "Control network access",
}

var aclAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add ACL rule",
	Long:  "add ACL rule\n\nExample:\n  cq acl add --groupid sg-fd8cc1ee --way ingress --protocol tcp --port 22 --address 192.0.2.0/24\n",
	Run: func(cmd *cobra.Command, args []string) {
		acl.Add(args)
	},
}

var aclDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete ACL rule",
	Long:  "delete ACL rule\n\nExample:\n  cq acl delete --groupid sg-fd8cc1ee --way ingress --protocol tcp --port 22 --address 192.0.2.0/24\n",
	Run: func(cmd *cobra.Command, args []string) {
		acl.Delete(args)
	},
}

var aclDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy ACL (CAN NOT RESTORE)",
	Long:  "destroy ACL (CAN NOT RESTORE)\n\nExample:\n  cq acl destroy sg-fd8cc1ee\n",
	Run: func(cmd *cobra.Command, args []string) {
		acl.Destroy(args)
	},
}

var ruleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List acl",
	Long:  "List acl",
	Run: func(cmd *cobra.Command, args []string) {
		acl.AclList()
	},
}

var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Show ACL rule",
	Long:  "Show ACL rule",
	Run: func(cmd *cobra.Command, args []string) {
		acl.RuleList(args)
	},
}
