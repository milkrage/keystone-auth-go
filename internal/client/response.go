package client

import "time"

type Response struct {
	Token Token `json:"token"`
}

type Token struct {
	ExpiresAt time.Time `json:"expires_at"`
}
