package domain

import "github.com/golang-jwt/jwt/v5"

type ServiceName string

const (
    ServiceCluster ServiceName = "cluster-service"
    ServiceRoute ServiceName   = "route-service"
)

type ServiceClaims struct {
	Service ServiceName		`json:"service"`
	jwt.RegisteredClaims
}