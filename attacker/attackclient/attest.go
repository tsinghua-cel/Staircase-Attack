package attackclient

import (
	"context"
	"github.com/tsinghua-cel/attacker-service/types"
)

var attestModule = "attest"

func (ec *Client) AttestBeforeBroadCast(ctx context.Context, slot uint64) (types.AttackerResponse, error) {
	var result types.AttackerResponse
	err := ec.c.CallContext(ctx, &result, attestModule+"_beforeBroadCast", slot)
	//err := ec.c.CallContext(ctx, &result, attestModule+"_beforeBroadCast", ec.clientInfo(), slot)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (ec *Client) AttestAfterBroadCast(ctx context.Context, slot uint64) (types.AttackerResponse, error) {
	var result types.AttackerResponse
	err := ec.c.CallContext(ctx, &result, attestModule+"_afterBroadCast", slot)
	//err := ec.c.CallContext(ctx, &result, attestModule+"_afterBroadCast", ec.clientInfo(), slot)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (ec *Client) AttestBeforeSign(ctx context.Context, slot uint64, pubkey string, attestDataBase64 string) (types.AttackerResponse, error) {
	var result types.AttackerResponse
	err := ec.c.CallContext(ctx, &result, attestModule+"_beforeSign", slot, pubkey, attestDataBase64)
	//err := ec.c.CallContext(ctx, &result, attestModule+"_beforeSign", ec.clientInfo(), slot, pubkey, attestDataBase64)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (ec *Client) AttestAfterSign(ctx context.Context, slot uint64, pubkey string, siginedAttestDataBase64 string) (types.AttackerResponse, error) {
	var result types.AttackerResponse
	err := ec.c.CallContext(ctx, &result, attestModule+"_afterSign", slot, pubkey, siginedAttestDataBase64)
	//err := ec.c.CallContext(ctx, &result, attestModule+"_afterSign", ec.clientInfo(), slot, pubkey, siginedAttestDataBase64)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (ec *Client) AttestBeforePropose(ctx context.Context, slot uint64, pubkey string, siginedAttestDataBase64 string) (types.AttackerResponse, error) {
	var result types.AttackerResponse
	err := ec.c.CallContext(ctx, &result, attestModule+"_beforePropose", slot, pubkey, siginedAttestDataBase64)
	//err := ec.c.CallContext(ctx, &result, attestModule+"_beforePropose", ec.clientInfo(), slot, pubkey, siginedAttestDataBase64)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (ec *Client) AttestAfterPropose(ctx context.Context, slot uint64, pubkey string, siginedAttestDataBase64 string) (types.AttackerResponse, error) {
	var result types.AttackerResponse
	err := ec.c.CallContext(ctx, &result, attestModule+"_afterPropose", slot, pubkey, siginedAttestDataBase64)
	//err := ec.c.CallContext(ctx, &result, attestModule+"_afterPropose", ec.clientInfo(), slot, pubkey, siginedAttestDataBase64)
	if err != nil {
		return result, err
	}
	return result, nil
}
