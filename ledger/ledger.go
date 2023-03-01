package ledger

import (
	"context"
	"time"

	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"

	"github.com/tron-us/go-btfs-common/crypto"
	"github.com/tron-us/go-btfs-common/utils/grpc"
	"github.com/tron-us/protobuf/proto"

	ic "github.com/libp2p/go-libp2p/core/crypto"
)

const LedgerVersion = "BTFS_Escrow_1.0.0"

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
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

func NewSignedChannelCommit(commit *ledgerpb.ChannelCommit, sig []byte) *ledgerpb.SignedChannelCommit {
	return &ledgerpb.SignedChannelCommit{
		Channel:   commit,
		Signature: sig,
	}
}

func NewSignedCreateAccountRequest(key *ledgerpb.PublicKey, sig []byte) *ledgerpb.SignedCreateAccountRequest {
	return &ledgerpb.SignedCreateAccountRequest{
		Key:                 key,
		Signature:           sig,
		ClientVersionNumber: LedgerVersion,
	}
}

func (c *Client) ImportAccount(ctx context.Context, pubKey ic.PubKey) (*ledgerpb.Account, error) {
	keyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	var res *ledgerpb.CreateAccountResult
	err = grpc.LedgerClient(c.addr).WithContext(ctx, func(ctx context.Context, client ledgerpb.ChannelsClient) error {
		res, err = client.CreateAccount(ctx, &ledgerpb.PublicKey{Key: keyBytes})
		return err
	})
	if err != nil {
		return nil, err
	}
	return res.GetAccount(), nil
}

func (c *Client) ImportSignedAccount(ctx context.Context, privKey ic.PrivKey, pubKey ic.PubKey) (*ledgerpb.SignedCreateAccountResult, error) {
	pubKeyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	signedPubKey := &ledgerpb.PublicKey{Key: pubKeyBytes}
	sigBytes, err := proto.Marshal(signedPubKey)
	if err != nil {
		return nil, err
	}
	signature, err := privKey.Sign(sigBytes)
	if err != nil {
		return nil, err
	}
	signedPubkey := NewSignedCreateAccountRequest(signedPubKey, signature)

	var result *ledgerpb.SignedCreateAccountResult
	err = grpc.LedgerClient(c.addr).WithContext(ctx, func(ctx context.Context, client ledgerpb.ChannelsClient) error {
		result, err = client.SignedCreateAccount(ctx, signedPubkey)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateChannel(ctx context.Context, channelCommit *ledgerpb.ChannelCommit, sig []byte) (*ledgerpb.ChannelID, error) {
	var (
		channelId *ledgerpb.ChannelID
		err       error
	)
	err = grpc.LedgerClient(c.addr).WithContext(ctx, func(ctx context.Context, client ledgerpb.ChannelsClient) error {
		channelId, err = client.CreateChannel(ctx, &ledgerpb.SignedChannelCommit{
			Channel:   channelCommit,
			Signature: sig,
		})
		return err
	})
	return channelId, err
}

func (c *Client) CloseChannel(ctx context.Context, signedChannelState *ledgerpb.SignedChannelState) error {
	return grpc.LedgerClient(c.addr).WithContext(ctx, func(ctx context.Context, client ledgerpb.ChannelsClient) error {
		_, err := client.CloseChannel(ctx, signedChannelState)
		return err
	})
}

func NewSignedPublicKey(privK ic.PrivKey, pubK ic.PubKey) (*ledgerpb.SignedPublicKey, error) {
	raw, err := ic.MarshalPublicKey(pubK)
	if err != nil {
		return nil, err
	}
	lgPubKey := &ledgerpb.PublicKey{
		Key: raw,
	}
	sig, err := crypto.Sign(privK, lgPubKey)
	if err != nil {
		return nil, err
	}
	return &ledgerpb.SignedPublicKey{
		Key:       lgPubKey,
		Signature: sig,
	}, nil
}
