package grpcauth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aclgo/simple-api-gateway/config"
	jwt "github.com/golang-jwt/jwt/v5"
)

type grpcAuth struct {
	mu         sync.Mutex
	token      string
	expiresAt  time.Time
	privateKey *rsa.PrivateKey
}

func NewGrpcAuth(cfg *config.Config) *grpcAuth {

	privKeyData, err := os.ReadFile(cfg.PathPrivatePem)
	if err != nil {
		log.Fatalf("os.ReadFile: %v", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		log.Fatalf("jwt.ParseEdPrivateKeyFromPEM: %v", err)
	}

	g := grpcAuth{
		mu:         sync.Mutex{},
		privateKey: privKey,
	}

	return &g
}

func (a *grpcAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.token == "" || time.Now().Add(10*time.Second).After(a.expiresAt) {
		newTtk, err := a.GenToken(ctx)
		if err != nil {
			return nil, err
		}

		a.token = newTtk
		a.expiresAt = time.Now().Add(15 * time.Minute)
	}

	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", a.token),
	}, nil
}

func (a *grpcAuth) RequireTransportSecurity() bool {
	return true
}

func (a *grpcAuth) GenToken(context.Context) (string, error) {
	claims := jwt.MapClaims{
		"iss": "simple-api-gateway",
		"sub": "gateway-client",
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	ttk := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ttkStr, err := ttk.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("falha at wrtie token rs256: %w", err)
	}

	return ttkStr, nil
}
