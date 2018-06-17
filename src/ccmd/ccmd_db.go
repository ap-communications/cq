package ccmd

import (
	"github.com/ap-communications/cq/src/ccmd/commons"
	"github.com/ap-communications/cq/src/ccmd/db"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dbCmd)

	dbCmd.AddCommand(dbDestroyCmd)
	dbDestroyCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "Destroy without confirmation")

	dbCmd.AddCommand(dbInspectCmd)
	dbInspectCmd.SetUsageTemplate("Usage:\n  cq db inspect\n    or\n  cq db inspect [DB-id] [DB-id] ...\n")

	dbCmd.AddCommand(dbListCmd)
	dbListCmd.Flags().StringVarP(&commons.Flags.Delimiter, "delimiter", "d", "\t", "delimiter")

	dbCmd.AddCommand(dbRebootCmd)
	dbRebootCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "Delete without confirmation")
	dbRebootCmd.Flags().BoolVarP(&commons.Flags.NoFailover, "no-failover", "", false, "execute to no failover reboot (it will be take more DB downtime)")

	dbCmd.AddCommand(dbSnapDelCmd)
	dbSnapDelCmd.Flags().BoolVarP(&commons.Flags.Force, "force", "f", false, "Delete without confirmation")

	dbCmd.AddCommand(dbSnapListCmd)
	dbSnapListCmd.Flags().StringVarP(&commons.Flags.Delimiter, "delimiter", "d", "\t", "delimiter")

	dbCmd.AddCommand(dbSnapshotCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Control manged DBs",
	Long:  "Control manged DBs",
}

var dbDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "[DANGER]  Destroy DB  (CAN NOT RESTORE)",
	Long:  "[DANGER]  Destroy DB  (CAN NOT RESTORE)",
	Run: func(cmd *cobra.Command, args []string) {
		db.Destroy(args)
	},
}

var dbInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Print detailed DB instance info",
	Long:  "Print detailed DB instance info",
	Run: func(cmd *cobra.Command, args []string) {
		db.Inspect(args)
	},
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "DB list of all regions",
	Long:  "DB list of all regions",
	Run: func(cmd *cobra.Command, args []string) {
		db.List()
	},
}

var dbRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot DB with failover (default)",
	Long:  "Reboot DB with failover (default)\nIf you want to no failover reboot, must set failover flag at false\n\n\n*** WARNING ***\n  If you execute to no failover reboot, it will be take more DB down time\n\n\nINFO:\n  If you want not configured MultiAZ or High Availability to DB instance, you must set --no-failover flag\n\n\n\n...\n",
	Run: func(cmd *cobra.Command, args []string) {
		db.Reboot(args)
	},
}

var dbSnapDelCmd = &cobra.Command{
	Use:   "snapdel",
	Short: "Delete DB snapshot",
	Long:  "Delete DB snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		db.SnapDel(args)
	},
}

var dbSnapListCmd = &cobra.Command{
	Use:   "snaplist",
	Short: "DB snapshot list of all regions",
	Long:  "DB snapshot list of all regions",
	Run: func(cmd *cobra.Command, args []string) {
		db.SnapList()
	},
}

var dbSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Take DB snapshot",
	Long:  "Take DB snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		db.Snapshot(args)
	},
}
