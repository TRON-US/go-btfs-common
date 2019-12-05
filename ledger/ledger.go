package ledger

import (
	"context"
	"time"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"
	"github.com/tron-us/protobuf/proto"
	"google.golang.org/grpc"
)

func NewClient(conn *grpc.ClientConn) ledgerpb.ChannelsClient {
	return ledgerpb.NewChannelsClient(conn)
}

func NewAccount(pubKey ic.PubKey, amount int64) (*ledgerpb.Account, error) {
	addr, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	return &ledgerpb.Account{
		Address: &ledgerpb.PublicKey{Key: addr},
		Balance: amount,
	}, nil
}

func NewChannelCommit(fromKey ic.PubKey, toKey ic.PubKey, amount int64) (*ledgerpb.ChannelCommit, error) {
	fromAddr, err := fromKey.Raw()
	if err != nil {
		return nil, err
	}
	toAddr, err := toKey.Raw()
	if err != nil {
		return nil, err
	}
	return &ledgerpb.ChannelCommit{
		Payer:     &ledgerpb.PublicKey{Key: fromAddr},
		Recipient: &ledgerpb.PublicKey{Key: toAddr},
		Amount:    amount,
		PayerId:   time.Now().UnixNano(),
	}, err
}

func NewChannelState(id *ledgerpb.ChannelID, sequence int64, fromAccount *ledgerpb.Account, toAccount *ledgerpb.Account) *ledgerpb.ChannelState {
	return &ledgerpb.ChannelState{
		Id:       id,
		Sequence: sequence,
		From:     fromAccount,
		To:       toAccount,
	}
}

func NewSignedChannelState(channelState *ledgerpb.ChannelState, fromSig []byte, toSig []byte) *ledgerpb.SignedChannelState {
	return &ledgerpb.SignedChannelState{
		Channel:       channelState,
		FromSignature: fromSig,
		ToSignature:   toSig,
	}
}

func ImportAccount(ctx context.Context, pubKey ic.PubKey, ledgerClient ledgerpb.ChannelsClient) (*ledgerpb.Account, error) {
	keyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	res, err := ledgerClient.CreateAccount(ctx, &ledgerpb.PublicKey{Key: keyBytes})
	if err != nil {
		return nil, err
	}
	return res.GetAccount(), nil
}

func ImportSignedAccount(ctx context.Context, privKey ic.PrivKey, pubKey ic.PubKey, ledgerClient ledgerpb.ChannelsClient) (*ledgerpb.SignedCreateAccountResult, error) {
	pubKeyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	singedPubKey := &ledgerpb.PublicKey{Key: pubKeyBytes}
	sigBytes, err := proto.Marshal(singedPubKey)
	signature, err := privKey.Sign(sigBytes)
	if err != nil {
		return nil, err
	}
	signedPubkey := &ledgerpb.SignedPublicKey{Key: singedPubKey, Signature: signature}
	return ledgerClient.SignedCreateAccount(ctx, signedPubkey)
}

func CreateChannel(ctx context.Context, ledgerClient ledgerpb.ChannelsClient, channelCommit *ledgerpb.ChannelCommit, sig []byte) (*ledgerpb.ChannelID, error) {
	return ledgerClient.CreateChannel(ctx, &ledgerpb.SignedChannelCommit{
		Channel:   channelCommit,
		Signature: sig,
	})
}

func CloseChannel(ctx context.Context, ledgerClient ledgerpb.ChannelsClient, signedChannelState *ledgerpb.SignedChannelState) error {
	_, err := ledgerClient.CloseChannel(ctx, signedChannelState)
	if err != nil {
		return err
	}
	return nil
}
