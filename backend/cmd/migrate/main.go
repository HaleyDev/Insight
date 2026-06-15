// Command migrate runs database schema migration against the configured
// PostgreSQL instance. It reuses the project's config loader and ORM
// manager, then invokes GORM AutoMigrate on every model struct.
//
// Usage:
//
//	go run ./cmd/migrate -c config -e local
//	go run ./cmd/migrate -c config -e local -seed
//	go run ./cmd/migrate -c config -e local -drop   # drop tables first
package main

import (
	"flag"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/pkg/config"
)

var (
	cfgDir = flag.String("c", "config", "config dir")
	env    = flag.String("e", "local", "env name, e.g. local / docker")
	seed   = flag.Bool("seed", false, "insert demo seed rows after creating tables")
	drop   = flag.Bool("drop", false, "drop existing tables before migrate (DANGEROUS)")
)

// allModels is the full list of GORM models managed by this migrator.
var allModels = []interface{}{
	&model.UserBaseModel{},
}

func main() {
	flag.Parse()

	_ = config.New(*cfgDir, config.WithEnv(*env))
	model.Init()

	db, err := model.GetDB()
	if err != nil {
		log.Fatalf("get default db failed: %v", err)
	}

	if *drop {
		log.Println("[drop] dropping existing tables ...")
		if err := db.Migrator().DropTable(allModels...); err != nil {
			log.Fatalf("drop table failed: %v", err)
		}
	}

	log.Println("[migrate] running AutoMigrate ...")
	if err := db.AutoMigrate(allModels...); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	log.Println("[migrate] OK: user_base")

	if *seed {
		seedDemoRows(db)
	}

	log.Println("migration finished.")
}

// seedDemoRows 插入一个默认管理员账号，密码：123456 (bcrypt cost=10)。
func seedDemoRows(db *gorm.DB) {
	log.Println("[seed] inserting default admin ...")

	now := time.Now()
	user := map[string]interface{}{
		"id":         1,
		"username":   "admin",
		"password":   "$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym",
		"email":      "admin@insight.com",
		"avatar":     model.DefaultAvatar,
		"role":       model.RoleAdmin,
		"created_at": now,
		"updated_at": now,
	}
	if err := db.Table("user_base").Create(user).Error; err != nil {
		log.Printf("[seed] user_base id=%v skip: %v", user["id"], err)
	}

	if err := db.Exec("SELECT setval(pg_get_serial_sequence('user_base', 'id'), COALESCE((SELECT MAX(id) FROM user_base), 1))").Error; err != nil {
		log.Printf("[seed] setval user_base skip: %v", err)
	}

	log.Println("[seed] done")
}
