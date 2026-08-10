package main

import (
	"fmt"
	"os"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/jwt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: service-token <cluster-service|route-service>")
		os.Exit(1)
	}

	serviceName := domain.ServiceName(os.Args[1])
	token, err := jwt.GenerateServiceToken(serviceName)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println(token)
}