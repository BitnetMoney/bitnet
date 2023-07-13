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

package rpc

import (
	"context"
	"net/http"
)

type mdHeaderKey struct{}

// NewContextWithHeaders wraps the given context, adding HTTP headers. These headers will
// be applied by Client when making a request using the returned context.
func NewContextWithHeaders(ctx context.Context, h http.Header) context.Context {
	if len(h) == 0 {
		// This check ensures the header map set in context will never be nil.
		return ctx
	}

	var ctxh http.Header
	prev, ok := ctx.Value(mdHeaderKey{}).(http.Header)
	if ok {
		ctxh = setHeaders(prev.Clone(), h)
	} else {
		ctxh = h.Clone()
	}
	return context.WithValue(ctx, mdHeaderKey{}, ctxh)
}

// headersFromContext is used to extract http.Header from context.
func headersFromContext(ctx context.Context) http.Header {
	source, _ := ctx.Value(mdHeaderKey{}).(http.Header)
	return source
}

// setHeaders sets all headers from src in dst.
func setHeaders(dst http.Header, src http.Header) http.Header {
	for key, values := range src {
		dst[http.CanonicalHeaderKey(key)] = values
	}
	return dst
}
