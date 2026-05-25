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
	&model.UserFansModel{},
	&model.UserFollowModel{},
	&model.UserStatModel{},
}

func main() {
	flag.Parse()

	// Bootstrap viper-based config so pkg/storage/orm can load
	// database.yaml from <cfgDir>/<env>/.
	_ = config.New(*cfgDir, config.WithEnv(*env))

	// Initialize ORM connections defined in database.yaml.
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
	log.Println("[migrate] OK: user_base, user_fans, user_follow, user_stat")

	if *seed {
		seedDemoRows(db)
	}

	log.Println("migration finished.")
}

// seedDemoRows inserts the demo rows shipped with the project. Uses
// ON CONFLICT DO NOTHING so the command is idempotent.
func seedDemoRows(db *gorm.DB) {
	log.Println("[seed] inserting demo rows ...")

	t1, _ := time.Parse("2006-01-02 15:04:05", "2020-05-23 00:12:30")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2020-05-23 00:23:10")

	users := []map[string]interface{}{
		{"id": 1, "username": "test-name", "password": "$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym", "avatar": "/uploads/avatar.jpg", "phone": int64(13010102020), "email": "123@cc.com", "sex": 1, "created_at": t1, "updated_at": t1},
		{"id": 2, "username": "admin", "password": "$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym", "avatar": "13010102021", "phone": int64(1), "email": "1234@cc.com", "sex": 0, "created_at": t1, "updated_at": t1},
		{"id": 3, "username": "admin2", "password": "$2a$10$Dps9oN3Oe3ZDMACih3DCGeTvR.jW/I8WD1NqapCJ6Vq3PzjnusI9i", "avatar": "13010102022", "phone": int64(0), "email": "12345@cc.com", "sex": 0, "created_at": t2, "updated_at": t2},
	}
	for _, u := range users {
		if err := db.Table("user_base").Create(u).Error; err != nil {
			log.Printf("[seed] user_base id=%v skip: %v", u["id"], err)
		}
	}

	follows := []map[string]interface{}{
		{"id": 1, "user_id": 1, "followed_uid": 2, "status": 1, "created_at": t1},
		{"id": 2, "user_id": 1, "followed_uid": 3, "status": 1, "created_at": t2},
	}
	for _, f := range follows {
		if err := db.Table("user_follow").Create(f).Error; err != nil {
			log.Printf("[seed] user_follow id=%v skip: %v", f["id"], err)
		}
	}

	fans := []map[string]interface{}{
		{"id": 1, "user_id": 2, "follower_uid": 1, "status": 1, "created_at": t1},
		{"id": 2, "user_id": 3, "follower_uid": 1, "status": 1, "created_at": t2},
	}
	for _, f := range fans {
		if err := db.Table("user_fans").Create(f).Error; err != nil {
			log.Printf("[seed] user_fans id=%v skip: %v", f["id"], err)
		}
	}

	stats := []map[string]interface{}{
		{"id": 1, "user_id": 1, "follow_count": 3, "follower_count": 0, "status": 1, "created_at": t1},
		{"id": 2, "user_id": 2, "follow_count": 0, "follower_count": 0, "status": 1, "created_at": t1},
		{"id": 8, "user_id": 3, "follow_count": 0, "follower_count": 1, "status": 1, "created_at": t2},
	}
	for _, s := range stats {
		if err := db.Table("user_stat").Create(s).Error; err != nil {
			log.Printf("[seed] user_stat id=%v skip: %v", s["id"], err)
		}
	}

	// Reset sequences to MAX(id) so subsequent INSERTs without explicit
	// id continue from the right value.
	for _, tbl := range []string{"user_base", "user_follow", "user_fans", "user_stat"} {
		if err := db.Exec("SELECT setval(pg_get_serial_sequence(?, 'id'), COALESCE((SELECT MAX(id) FROM "+tbl+"), 1))", tbl).Error; err != nil {
			log.Printf("[seed] setval %s skip: %v", tbl, err)
		}
	}

	log.Println("[seed] done")
}
