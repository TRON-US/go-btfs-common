package ledger

import (
	"fmt"

	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"

	ic "github.com/libp2p/go-libp2p/core/crypto"
)

func NewPayinRequest(payinId string, payerPubkey ic.PubKey, state *ledgerpb.SignedChannelState) (
	*escrowpb.PayinRequest, error) {
	raw, err := ic.MarshalPublicKey(payerPubkey)
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
	raw, err := ic.MarshalPublicKey(key)
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
	amount int64, ps escrowpb.Schedule, period int32, contrType escrowpb.ContractType,
	contingentAmount int64, storageLength int) (*escrowpb.EscrowContract, error) {
	payerAddr, err := ic.MarshalPublicKey(payerPubKey)
	if err != nil {
		return nil, err
	}
	var hostAddr []byte
	if hostPubKey != nil {
		hostAddr, err = ic.MarshalPublicKey(hostPubKey)
		if err != nil {
			return nil, err
		}
	}
	authAddress, err := ic.MarshalPublicKey(authPubKey)
	if err != nil {
		return nil, err
	}
	numPayouts := 1
	switch ps {
	case escrowpb.Schedule_MONTHLY:
		numPayouts = storageLength / 30
	case escrowpb.Schedule_QUARTERLY:
		numPayouts = storageLength / 30 / 3
	case escrowpb.Schedule_ANNUALLY:
		numPayouts = storageLength / 30 / 3 / 4
	case escrowpb.Schedule_CUSTOMIZED:
		numPayouts = storageLength / int(period)
	default:
		return nil, fmt.Errorf("invalide PayoutSchedule: %v", ps)
	}
	if numPayouts == 0 {
		numPayouts = 1
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
		NumPayouts:            int32(numPayouts),
		CustomizePayoutPeriod: period,
		Type:                  contrType,
		ContingentAmount:      contingentAmount,
	}, nil
}
