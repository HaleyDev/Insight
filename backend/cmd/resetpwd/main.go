package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/pkg/config"
)

// 一次性脚本：为 admin 重置密码为 123456
func main() {
	_ = config.New("config", config.WithEnv("local"))
	model.Init()

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt err: %v", err)
	}
	fmt.Printf("hash: %s\n", string(hash))

	db, err := model.GetDB()
	if err != nil {
		log.Fatalf("get db err: %v", err)
	}
	if err := db.Table("user_base").Where("email = ?", "admin@insight.com").
		Update("password", string(hash)).Error; err != nil {
		log.Fatalf("update err: %v", err)
	}
	fmt.Println("admin password reset to 123456")
}
