package keystoneauth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/milkrage/keystone-auth-go/internal/cache"
	"github.com/milkrage/keystone-auth-go/internal/client"
)

type Authenticator struct {
	credentials Credentials
	cache       cache.Cache
	client      client.Client
	mu          sync.Mutex
}

func NewAuthenticator(ctx context.Context, url string, credentials Credentials) (*Authenticator, error) {
	a := &Authenticator{
		credentials: credentials,
		client:      client.Client{BaseURL: url},
	}

	_, err := a.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Authenticator) GetToken(ctx context.Context) (string, error) {
	token, ok := a.cache.Get()
	if ok {
		return token, nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-Checked Locking Pattern (DCL).
	token, ok = a.cache.Get()
	if ok {
		return token, nil
	}

	token, expiresAt, err := a.client.IssueToken(ctx, a.credentials.toRequest())
	if err != nil {
		return "", fmt.Errorf("failed to issue token: %w", err)
	}
	a.cache.Set(token, expiresAt)

	return token, nil
}

func (a *Authenticator) InvalidateToken(t string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	current, _ := a.cache.Get()
	if current == t {
		a.cache.Set("", time.Time{})
	}
}
