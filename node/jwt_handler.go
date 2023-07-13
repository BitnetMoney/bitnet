// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package node

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtExpiryTimeout = 60 * time.Second

type jwtHandler struct {
	keyFunc func(token *jwt.Token) (interface{}, error)
	next    http.Handler
}

// newJWTHandler creates a http.Handler with jwt authentication support.
func newJWTHandler(secret []byte, next http.Handler) http.Handler {
	return &jwtHandler{
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		next: next,
	}
}

// ServeHTTP implements http.Handler
func (handler *jwtHandler) ServeHTTP(out http.ResponseWriter, r *http.Request) {
	var (
		strToken string
		claims   jwt.RegisteredClaims
	)
	if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		strToken = strings.TrimPrefix(auth, "Bearer ")
	}
	if len(strToken) == 0 {
		http.Error(out, "missing token", http.StatusUnauthorized)
		return
	}
	// We explicitly set only HS256 allowed, and also disables the
	// claim-check: the RegisteredClaims internally requires 'iat' to
	// be no later than 'now', but we allow for a bit of drift.
	token, err := jwt.ParseWithClaims(strToken, &claims, handler.keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithoutClaimsValidation())

	switch {
	case err != nil:
		http.Error(out, err.Error(), http.StatusUnauthorized)
	case !token.Valid:
		http.Error(out, "invalid token", http.StatusUnauthorized)
	case !claims.VerifyExpiresAt(time.Now(), false): // optional
		http.Error(out, "token is expired", http.StatusUnauthorized)
	case claims.IssuedAt == nil:
		http.Error(out, "missing issued-at", http.StatusUnauthorized)
	case time.Since(claims.IssuedAt.Time) > jwtExpiryTimeout:
		http.Error(out, "stale token", http.StatusUnauthorized)
	case time.Until(claims.IssuedAt.Time) > jwtExpiryTimeout:
		http.Error(out, "future token", http.StatusUnauthorized)
	default:
		handler.next.ServeHTTP(out, r)
	}
}
