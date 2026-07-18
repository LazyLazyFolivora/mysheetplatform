package pkg

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type codeEntry struct {
	Code      string
	ExpiresAt time.Time
}

var (
	codeStore = sync.Map{}
)

func init() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			now := time.Now()
			codeStore.Range(func(key, value any) bool {
				if value.(codeEntry).ExpiresAt.Before(now) {
					codeStore.Delete(key)
				}
				return true
			})
		}
	}()
}

func GenerateAndStoreCode(email string) string {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	codeStore.Store(email, codeEntry{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	return code
}

func VerifyCode(email, code string) bool {
	val, ok := codeStore.Load(email)
	if !ok {
		return false
	}
	entry := val.(codeEntry)
	if time.Now().After(entry.ExpiresAt) {
		codeStore.Delete(email)
		return false
	}
	if entry.Code != code {
		return false
	}
	codeStore.Delete(email)
	return true
}
