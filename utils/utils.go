package utils

import (
	"runtime"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
