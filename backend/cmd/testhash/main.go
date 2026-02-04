package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	
	// 生成新的 hash
	newHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Printf("admin123 的正确 hash: %s\n", string(newHash))
	
	// 验证
	err := bcrypt.CompareHashAndPassword(newHash, []byte(password))
	fmt.Printf("验证结果: %v (true=正确)\n", err == nil)
}

