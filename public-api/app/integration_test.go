package app

import (
	"net/http"
	"testing"
	"time"

	usrsvc "github.com/vaskoz/monorepo-multi-service/user-service/app"
)

func TestPublicAPIUsers(t *testing.T) {
	go usrsvc.RealMain()
	go RealMain()

	time.Sleep(1 * time.Second)

	publicApiClient := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/users", nil)
	publicApiClient.Do(req)
}

func BenchmarkPublicAPIUsers(b *testing.B) {
	publicApiClient := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/users", nil)

	for i := 0; i < b.N; i++ {
		publicApiClient.Do(req)
	}
}
