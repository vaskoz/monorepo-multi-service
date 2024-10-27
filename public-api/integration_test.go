package main

import (
	"net/http"
	"testing"

	usrsvc "github.com/vaskoz/monorepo-multi-service/user-service"
)

func TestPublicAPIUsers(t *testing.T) {
	usrsvc.StartMain()
	go main()

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
