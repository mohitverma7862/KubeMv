package auth

import "context"

// Principal is the signed-in operator. It carries no Kubernetes credentials.
type Principal struct {
	Username    string
	DisplayName string
	Roles       []Role
}

// Authenticator verifies an operator. Implementations must not retain the
// plaintext password after Authenticate returns.
type Authenticator interface {
	Authenticate(ctx context.Context, username, password string) (Principal, error)
}
