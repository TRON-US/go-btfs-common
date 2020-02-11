package ledger

import (
	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"

	ic "github.com/libp2p/go-libp2p-core/crypto"
)

func NewPayinRequest(payinId string, payerPubkey ic.PubKey, state *ledgerpb.SignedChannelState) (
	*escrowpb.PayinRequest, error) {
	raw, err := ic.RawFull(payerPubkey)
	if err != nil {
		return nil, err
	}
	return &escrowpb.PayinRequest{
		PayinId:           payinId,
		BuyerAddress:      raw,
		BuyerChannelState: state,
	}, nil
}

func NewSignedPayinRequest(req *escrowpb.PayinRequest, sig []byte) *escrowpb.SignedPayinRequest {
	return &escrowpb.SignedPayinRequest{
		Request:        req,
		BuyerSignature: sig,
	}
}

func NewContractID(id string, key ic.PubKey) (*escrowpb.ContractID, error) {
	raw, err := ic.RawFull(key)
	if err != nil {
		return nil, err
	}
	return &escrowpb.ContractID{
		ContractId: id,
		Address:    raw,
	}, nil
}

func NewSingedContractID(id *escrowpb.ContractID, sig []byte) *escrowpb.SignedContractID {
	return &escrowpb.SignedContractID{
		Data:      id,
		Signature: sig,
	}
}

func NewEscrowContract(id string, payerPubKey ic.PubKey, hostPubKey ic.PubKey, authPubKey ic.PubKey,
	amount int64, ps escrowpb.Schedule, period int32) (*escrowpb.EscrowContract, error) {
	payerAddr, err := ic.RawFull(payerPubKey)
	if err != nil {
		return nil, err
	}
	hostAddr, err := ic.RawFull(hostPubKey)
	if err != nil {
		return nil, err
	}
	authAddress, err := ic.RawFull(authPubKey)
	if err != nil {
		return nil, err
	}
	return &escrowpb.EscrowContract{
		ContractId:            id,
		BuyerAddress:          payerAddr,
		SellerAddress:         hostAddr,
		AuthAddress:           authAddress,
		Amount:                amount,
		CollateralAmount:      0,
		WithholdAmount:        0,
		TokenType:             escrowpb.TokenType_BTT,
		PayoutSchedule:        ps,
		NumPayouts:            1,
		CustomizePayoutPeriod: period,
	}, nil
}
