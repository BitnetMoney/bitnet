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

package core

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type AuditLogger struct {
	log log.Logger
	api ExternalAPI
}

func (l *AuditLogger) List(ctx context.Context) ([]common.Address, error) {
	l.log.Info("List", "type", "request", "metadata", MetadataFromContext(ctx).String())
	res, e := l.api.List(ctx)
	l.log.Info("List", "type", "response", "data", res)

	return res, e
}

func (l *AuditLogger) New(ctx context.Context) (common.Address, error) {
	return l.api.New(ctx)
}

func (l *AuditLogger) SignTransaction(ctx context.Context, args apitypes.SendTxArgs, methodSelector *string) (*ethapi.SignTransactionResult, error) {
	sel := "<nil>"
	if methodSelector != nil {
		sel = *methodSelector
	}
	l.log.Info("SignTransaction", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"tx", args.String(),
		"methodSelector", sel)

	res, e := l.api.SignTransaction(ctx, args, methodSelector)
	if res != nil {
		l.log.Info("SignTransaction", "type", "response", "data", common.Bytes2Hex(res.Raw), "error", e)
	} else {
		l.log.Info("SignTransaction", "type", "response", "data", res, "error", e)
	}
	return res, e
}

func (l *AuditLogger) SignData(ctx context.Context, contentType string, addr common.MixedcaseAddress, data interface{}) (hexutil.Bytes, error) {
	marshalledData, _ := json.Marshal(data) // can ignore error, marshalling what we just unmarshalled
	l.log.Info("SignData", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"addr", addr.String(), "data", marshalledData, "content-type", contentType)
	b, e := l.api.SignData(ctx, contentType, addr, data)
	l.log.Info("SignData", "type", "response", "data", common.Bytes2Hex(b), "error", e)
	return b, e
}

func (l *AuditLogger) SignGnosisSafeTx(ctx context.Context, addr common.MixedcaseAddress, gnosisTx GnosisSafeTx, methodSelector *string) (*GnosisSafeTx, error) {
	sel := "<nil>"
	if methodSelector != nil {
		sel = *methodSelector
	}
	data, _ := json.Marshal(gnosisTx) // can ignore error, marshalling what we just unmarshalled
	l.log.Info("SignGnosisSafeTx", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"addr", addr.String(), "data", string(data), "selector", sel)
	res, e := l.api.SignGnosisSafeTx(ctx, addr, gnosisTx, methodSelector)
	if res != nil {
		data, _ := json.Marshal(res) // can ignore error, marshalling what we just unmarshalled
		l.log.Info("SignGnosisSafeTx", "type", "response", "data", string(data), "error", e)
	} else {
		l.log.Info("SignGnosisSafeTx", "type", "response", "data", res, "error", e)
	}
	return res, e
}

func (l *AuditLogger) SignTypedData(ctx context.Context, addr common.MixedcaseAddress, data apitypes.TypedData) (hexutil.Bytes, error) {
	l.log.Info("SignTypedData", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"addr", addr.String(), "data", data)
	b, e := l.api.SignTypedData(ctx, addr, data)
	l.log.Info("SignTypedData", "type", "response", "data", common.Bytes2Hex(b), "error", e)
	return b, e
}

func (l *AuditLogger) EcRecover(ctx context.Context, data hexutil.Bytes, sig hexutil.Bytes) (common.Address, error) {
	l.log.Info("EcRecover", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"data", common.Bytes2Hex(data), "sig", common.Bytes2Hex(sig))
	b, e := l.api.EcRecover(ctx, data, sig)
	l.log.Info("EcRecover", "type", "response", "address", b.String(), "error", e)
	return b, e
}

func (l *AuditLogger) Version(ctx context.Context) (string, error) {
	l.log.Info("Version", "type", "request", "metadata", MetadataFromContext(ctx).String())
	data, err := l.api.Version(ctx)
	l.log.Info("Version", "type", "response", "data", data, "error", err)
	return data, err
}

func NewAuditLogger(path string, api ExternalAPI) (*AuditLogger, error) {
	l := log.New("api", "signer")
	handler, err := log.FileHandler(path, log.LogfmtFormat())
	if err != nil {
		return nil, err
	}
	l.SetHandler(handler)
	l.Info("Configured", "audit log", path)
	return &AuditLogger{l, api}, nil
}
