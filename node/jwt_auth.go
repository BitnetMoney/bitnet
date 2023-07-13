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
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang-jwt/jwt/v4"
)

// NewJWTAuth creates an rpc client authentication provider that uses JWT. The
// secret MUST be 32 bytes (256 bits) as defined by the Engine-API authentication spec.
//
// See https://github.com/ethereum/execution-apis/blob/main/src/engine/authentication.md
// for more details about this authentication scheme.
func NewJWTAuth(jwtsecret [32]byte) rpc.HTTPAuth {
	return func(h http.Header) error {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iat": &jwt.NumericDate{Time: time.Now()},
		})
		s, err := token.SignedString(jwtsecret[:])
		if err != nil {
			return fmt.Errorf("failed to create JWT token: %w", err)
		}
		h.Set("Authorization", "Bearer "+s)
		return nil
	}
}
