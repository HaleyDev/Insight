package migrate

import (
	"insight/data"
	"insight/internal/model"
	log "insight/internal/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "migrate",
		Short:   "Database migration tool",
		Example: "insight migrate",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initialize database connection
			data.InitData()
		},
		Run: func(cmd *cobra.Command, args []string) {
			migrate()
		},
	}
)

func migrate() {
	log.Logger.Info("Starting database migration...")

	// Auto migrate table schema
	err := data.MysqlDB.AutoMigrate(
		&model.AdminUser{},
		&model.Permission{},
	)

	if err != nil {
		log.Logger.Error("Database migration failed: " + err.Error())
		return
	}

	log.Logger.Info("Database migration completed!")
	log.Logger.Info("Created/Updated tables:")
	log.Logger.Info("  - a_admin_user (Admin user table)")
	log.Logger.Info("  - permissions (Permission table)")

	log.Logger.Info("Database migration completed")
}
