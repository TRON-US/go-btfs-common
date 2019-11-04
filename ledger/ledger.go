package ledger

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	ledgerPb "github.com/tron-us/go-btfs-common/protos/ledger"
	"github.com/tron-us/go-common/log"

	"github.com/gogo/protobuf/proto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"go.uber.org/zap"
)


func LedgerConnection(ledgerAddr, certFile string) (*grpc.ClientConn, error) {
	var grpcOption grpc.DialOption
	if certFile == "" {
		grpcOption = grpc.WithInsecure()
	} else {
		b := []byte(certFile)
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}
		credential := credentials.NewTLS(&tls.Config{RootCAs: cp})
		grpcOption = grpc.WithTransportCredentials(credential)
	}
	conn, err := grpc.Dial(ledgerAddr, grpcOption)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CloseConnection(conn *grpc.ClientConn) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Error("Failed to close connection: ", zap.Error(err))
		}
	}
}

func NewClient(conn *grpc.ClientConn) ledgerPb.ChannelsClient {
	return ledgerPb.NewChannelsClient(conn)
}

func NewAccount(pubKey ic.PubKey, amount int64) (*ledgerPb.Account, error) {
	addr, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	return &ledgerPb.Account{
		Address: &ledgerPb.PublicKey{Key: addr},
		Balance: amount,
	}, nil
}

func NewChannelCommit(fromKey ic.PubKey, toKey ic.PubKey, amount int64) (*ledgerPb.ChannelCommit, error) {
	fromAddr, err := fromKey.Raw()
	if err != nil {
		return nil, err
	}
	toAddr, err := toKey.Raw()
	if err != nil {
		return nil, err
	}
	return &ledgerPb.ChannelCommit{
		Payer:     &ledgerPb.PublicKey{Key: fromAddr},
		Recipient: &ledgerPb.PublicKey{Key: toAddr},
		Amount:    amount,
		PayerId:   time.Now().UnixNano(),
	}, err
}

func NewChannelState(id *ledgerPb.ChannelID, sequence int64, fromAccount *ledgerPb.Account, toAccount *ledgerPb.Account) *ledgerPb.ChannelState {
	return &ledgerPb.ChannelState{
		Id:       id,
		Sequence: sequence,
		From:     fromAccount,
		To:       toAccount,
	}
}

func NewSignedChannelState(channelState *ledgerPb.ChannelState, fromSig []byte, toSig []byte) *ledgerPb.SignedChannelState {
	return &ledgerPb.SignedChannelState{
		Channel:       channelState,
		FromSignature: fromSig,
		ToSignature:   toSig,
	}
}

func ImportAccount(ctx context.Context, pubKey ic.PubKey, ledgerClient ledgerPb.ChannelsClient) (*ledgerPb.Account, error) {
	keyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	res, err := ledgerClient.CreateAccount(ctx, &ledgerPb.PublicKey{Key: keyBytes})
	if err != nil {
		return nil, err
	}
	return res.GetAccount(), nil
}

func ImportSignedAccount(ctx context.Context, privKey ic.PrivKey, pubKey ic.PubKey, ledgerClient ledgerPb.ChannelsClient) (*ledgerPb.SignedCreateAccountResult, error) {
	pubKeyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	singedPubKey := &ledgerPb.PublicKey{Key: pubKeyBytes}
	sigBytes, err := proto.Marshal(singedPubKey)
	signature, err := privKey.Sign(sigBytes)
	if err != nil {
		return nil, err
	}
	signedPubkey := &ledgerPb.SignedPublicKey{Key: singedPubKey, Signature: signature}
	return ledgerClient.SignedCreateAccount(ctx, signedPubkey)
}

func CreateChannel(ctx context.Context, ledgerClient ledgerPb.ChannelsClient, channelCommit *ledgerPb.ChannelCommit, sig []byte) (*ledgerPb.ChannelID, error) {
	return ledgerClient.CreateChannel(ctx, &ledgerPb.SignedChannelCommit{
		Channel:   channelCommit,
		Signature: sig,
	})
}

func CloseChannel(ctx context.Context, ledgerClient ledgerPb.ChannelsClient, signedChannelState *ledgerPb.SignedChannelState) error {
	_, err = ledgerClient.CloseChannel(ctx, signedChannelState)
	if err != nil {
		return err
	}
	return nil
}

func Sign(key ic.PrivKey, channelMessage proto.Message) ([]byte, error) {
	raw, err := proto.Marshal(channelMessage)
	if err != nil {
		return nil, err
	}
	return key.Sign(raw)
}

func Verify(key ic.PubKey, channelMessage proto.Message, sig []byte) (bool, error) {
	raw, err := proto.Marshal(channelMessage)
	if err != nil {
		return false, err
	}
	return key.Verify(raw, sig)
}
