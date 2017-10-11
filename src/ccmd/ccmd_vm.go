package ccmd

import (
	"ccmd/commons"
	"ccmd/vm"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(vmCmd)

	vmCmd.AddCommand(vmListCmd)
	vmListCmd.Flags().StringVarP(&commons.Flags.Delimiter, "delimiter", "d", "\t", "delimiter")

	vmCmd.AddCommand(vmStartCmd)
	vmStartCmd.SetUsageTemplate("Usage:\n  cq vm start [instance-id] [instance-id] ...\n")

	vmCmd.AddCommand(vmStopCmd)
	vmStopCmd.SetUsageTemplate("Usage:\n  cq vm stop [instance-id] [instance-id] ...\n")
	vmStopCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "without confirmation")

	vmCmd.AddCommand(vmRebootCmd)
	vmRebootCmd.SetUsageTemplate("Usage:\n  cq vm reboot [instance-id] [instance-id] ...\n")
	vmRebootCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "without confirmation")

	vmCmd.AddCommand(vmDestroyCmd)
	vmDestroyCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "without confirmation")

	vmCmd.AddCommand(vmBackupCmd)
	vmBackupCmd.SetUsageTemplate("Usage:  cq vm backup [instance-id] [instance-id] ...\n")

	vmCmd.AddCommand(vmRestoreCmd)
	vmRestoreCmd.SetUsageTemplate("Usage:\n  cq vm restore [filename] [filename] ...")

	vmCmd.AddCommand(vmEasyupCmd)
	vmEasyupCmd.Flags().StringVarP(&commons.Flags.Region, "region", "", commons.DEFAULT_REGION, "Region")
	vmEasyupCmd.Flags().StringVarP(&commons.Flags.ImageId, "imageid", "", "", "Instance image ID (default latest Amazon Linux image)")
	vmEasyupCmd.Flags().StringVarP(&commons.Flags.Type, "type", "", "t2.micro", "Instance type")
	vmEasyupCmd.Flags().StringVarP(&commons.Flags.Keyname, "key", "", "", "SSH key-pair (default new generate)")
	vmEasyupCmd.Flags().StringVarP(&commons.Flags.GroupId, "groupid", "", "", "Security Group ID (default new generate)")

	vmCmd.AddCommand(vmInspectCmd)
	vmInspectCmd.SetUsageTemplate("Usage:\n  cq vm inspect\n    or\n  cq vm inspect [instance-id] [instance-id] ...")
}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Control virtual machines",
	Long:  "Control virtual machines",
}

var vmListCmd = &cobra.Command{
	Use:   "list",
	Short: "VM list of all regions",
	Long:  "VM list of all regions",
	Run: func(cmd *cobra.Command, args []string) {
		vm.List()
	},
}

var vmStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Startup VM",
	Long:  "Startup VM",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Start(args)
	},
}

var vmStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Shutdown VM",
	Long:  "Shutdown VM",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Stop(args)
	},
}

var vmRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot running status VM",
	Long:  "Reboot running status VM",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Reboot(args)
	},
}
var vmDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "[DANGER]  Destroy VM  (CAN NOT RESTORE)",
	Long:  "[DANGER]  Destroy VM  (CAN NOT RESTORE)",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Destroy(args)
	},
}

var vmBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup configure information",
	Long:  "Backup configure information",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Backup(args)
	},
}

var vmRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore saved configure information file (.cq file)",
	Long:  "Restore saved configure information file (.cq file)",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Restore(args)
	},
}

var vmEasyupCmd = &cobra.Command{
	Use:   "easyup",
	Short: "Create and run new VM",
	Long:  "Create and run new VM",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Easyup()
	},
}

var vmInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Print detailed instance info",
	Long:  "Print detailed instance info",
	Run: func(cmd *cobra.Command, args []string) {
		vm.Inspect(args)
	},
}
