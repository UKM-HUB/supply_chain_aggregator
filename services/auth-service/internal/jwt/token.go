// Package jwt re-exports the shared JWT manager from supply-chain-aggregator/pkg/jwt.
package jwt

import pkgjwt "supply-chain-aggregator/pkg/jwt"

// CustomClaims is the shared claims type.
type CustomClaims = pkgjwt.CustomClaims

// Manager signs and validates JWTs. It delegates to pkg/jwt.Manager.
type Manager = pkgjwt.Manager

// NewManager creates a new JWT Manager with a 24-hour TTL.
func NewManager(secret string) *Manager {
	return pkgjwt.NewManager(secret, 0)
}
