package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "12345678"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Printf("UPDATE administrateurs SET password = '%s' WHERE email = 'johannliebert@gmail.com';\n", hash)
}