package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"

	hash1 := "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi"
	hash2 := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	err1 := bcrypt.CompareHashAndPassword([]byte(hash1), []byte(password))
	err2 := bcrypt.CompareHashAndPassword([]byte(hash2), []byte(password))

	fmt.Printf("Hash1 验证 'admin123': %v\n", err1 == nil)
	fmt.Printf("Hash2 验证 'admin123': %v\n", err2 == nil)

	// 生成一个新的正确哈希
	newHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Printf("\n新生成的 admin123 哈希:\n%s\n", string(newHash))
}
