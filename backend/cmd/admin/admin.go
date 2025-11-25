package admin

import (
	"fmt"
	"insight/data"
	"insight/internal/model"
	log "insight/internal/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "admin",
		Short:   "Admin user management tool",
		Example: "insight admin create --username=admin --password=123456",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initialize database connection
			data.InitData()
		},
	}

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new admin user",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initialize database connection
			data.InitData()
		},
		Run: createAdmin,
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all admin users",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initialize database connection
			data.InitData()
		},
		Run: listAdmins,
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an admin user by username",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initialize database connection
			data.InitData()
		},
		Run: deleteAdmin,
	}

	resetPwdCmd = &cobra.Command{
		Use:   "reset-password",
		Short: "Reset admin user password",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initialize database connection
			data.InitData()
		},
		Run: resetPassword,
	}

	// Flags
	username    string
	password    string
	email       string
	mobile      string
	nickname    string
	isAdmin     bool
	targetUser  string
	newPassword string
)

func init() {
	// Add subcommands
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(resetPwdCmd)

	// Create command flags
	createCmd.Flags().StringVarP(&username, "username", "u", "", "Username (required)")
	createCmd.Flags().StringVarP(&password, "password", "p", "", "Password (required)")
	createCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	createCmd.Flags().StringVarP(&mobile, "mobile", "m", "", "Mobile number")
	createCmd.Flags().StringVarP(&nickname, "nickname", "n", "", "Nickname")
	createCmd.Flags().BoolVar(&isAdmin, "admin", true, "Is admin user")
	createCmd.MarkFlagRequired("username")
	createCmd.MarkFlagRequired("password")

	// Delete command flags
	deleteCmd.Flags().StringVarP(&targetUser, "username", "u", "", "Username to delete (required)")
	deleteCmd.MarkFlagRequired("username")

	// Reset password command flags
	resetPwdCmd.Flags().StringVarP(&targetUser, "username", "u", "", "Username (required)")
	resetPwdCmd.Flags().StringVarP(&newPassword, "password", "p", "", "New password (required)")
	resetPwdCmd.MarkFlagRequired("username")
	resetPwdCmd.MarkFlagRequired("password")
}

func createAdmin(cmd *cobra.Command, args []string) {
	log.Logger.Info("Creating admin user: " + username)

	// Check database connection
	if data.MysqlDB == nil {
		log.Logger.Error("Database connection not initialized")
		return
	}

	// Test database connection
	sqlDB, err := data.MysqlDB.DB()
	if err != nil {
		log.Logger.Error("Failed to get database instance: " + err.Error())
		return
	}

	if err := sqlDB.Ping(); err != nil {
		log.Logger.Error("Failed to ping database: " + err.Error())
		return
	}

	// Check if username already exists
	adminUser := model.NewAdminUsers()
	existingUser := adminUser.GetUserInfo(username)
	if existingUser != nil {
		log.Logger.Warn("Username already exists: " + username)
		return
	}

	// Set default nickname if not provided
	if nickname == "" {
		nickname = username
	}

	// Create new admin user
	newUser := &model.AdminUser{
		Username: username,
		Password: password,
		Email:    email,
		Mobile:   mobile,
		NickName: nickname,
		Status:   1, // Active
		IsAdmin:  1, // Admin user
	}

	if !isAdmin {
		newUser.IsAdmin = 0
	}

	// Hash password manually
	hashedPassword, hashErr := newUser.PasswordHash(password)
	if hashErr != nil {
		log.Logger.Error("Failed to hash password: " + hashErr.Error())
		return
	}

	// Set the hashed password
	newUser.Password = hashedPassword

	// Use GORM Create but with proper model handling
	result := data.MysqlDB.Create(newUser)
	if result.Error != nil {
		log.Logger.Error("Failed to create admin user: " + result.Error.Error())
		return
	}

	log.Logger.Info("Admin user created successfully: " + username)
	log.Logger.Info("Username: " + username)
	log.Logger.Info("Nickname: " + nickname)
	log.Logger.Info("Email: " + email)
	log.Logger.Info("Mobile: " + mobile)
	if isAdmin {
		log.Logger.Info("Is Admin: true")
	} else {
		log.Logger.Info("Is Admin: false")
	}
}

func listAdmins(cmd *cobra.Command, args []string) {
	log.Logger.Info("Admin Users List:")
	log.Logger.Info("==================")

	var users []model.AdminUser
	result := data.MysqlDB.Find(&users)
	if result.Error != nil {
		log.Logger.Error("Failed to fetch users: " + result.Error.Error())
		return
	}

	if len(users) == 0 {
		log.Logger.Info("No admin users found.")
		return
	}

	fmt.Printf("%-5s %-15s %-15s %-25s %-15s %-8s %-8s\n",
		"ID", "Username", "Nickname", "Email", "Mobile", "IsAdmin", "Status")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	for _, user := range users {
		status := "Active"
		if user.Status != 1 {
			status = "Inactive"
		}

		isAdminStr := "No"
		if user.IsAdmin == 1 {
			isAdminStr = "Yes"
		}

		fmt.Printf("%-5d %-15s %-15s %-25s %-15s %-8s %-8s\n",
			user.ID, user.Username, user.NickName, user.Email, user.Mobile, isAdminStr, status)
	}
}

func deleteAdmin(cmd *cobra.Command, args []string) {
	log.Logger.Info("Deleting admin user: " + targetUser)

	// Check if user exists
	adminUser := model.NewAdminUsers()
	user := adminUser.GetUserInfo(targetUser)
	if user == nil {
		log.Logger.Warn("User not found: " + targetUser)
		return
	}

	// Soft delete the user
	result := data.MysqlDB.Delete(user)
	if result.Error != nil {
		log.Logger.Error("Failed to delete admin user: " + result.Error.Error())
		return
	}

	log.Logger.Info("User deleted successfully: " + targetUser)
}

func resetPassword(cmd *cobra.Command, args []string) {
	log.Logger.Info("Resetting password for user: " + targetUser)

	// Check if user exists
	adminUser := model.NewAdminUsers()
	user := adminUser.GetUserInfo(targetUser)
	if user == nil {
		log.Logger.Warn("User not found: " + targetUser)
		return
	}

	// Update password
	user.Password = newPassword
	err := user.ChangePassword()
	if err != nil {
		log.Logger.Error("Failed to reset password: " + err.Error())
		return
	}

	log.Logger.Info("Password reset successfully for user: " + targetUser)
}
