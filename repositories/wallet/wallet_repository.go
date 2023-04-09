//go:generate go run github.com/golang/mock/mockgen -destination ../../mocks/$GOFILE -package mocks . WalletRepository

package walletrepository

import (
	"context"
	"net/http"

	walletApi "buf.build/gen/go/ride/wallet/bufbuild/connect-go/ride/wallet/v1alpha1/walletv1alpha1connect"
	pb "buf.build/gen/go/ride/wallet/protocolbuffers/go/ride/wallet/v1alpha1"
	"github.com/bufbuild/connect-go"
	"github.com/ride-app/driver-service/config"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, id string) (*pb.Wallet, error)
}

type Impl struct {
	walletApi walletApi.WalletServiceClient
}

func New() (*Impl, error) {
	client := walletApi.NewWalletServiceClient(
		http.DefaultClient,
		config.Env.Wallet_Service_Host,
	)

	return &Impl{walletApi: client}, nil
}

func (r *Impl) GetWallet(ctx context.Context, id string) (*pb.Wallet, error) {
	res, err := r.walletApi.GetWallet(ctx, connect.NewRequest(&pb.GetWalletRequest{
		Name: "users/" + id + "/wallet",
	}))

	if err != nil {
		return nil, err
	}

	return res.Msg.Wallet, nil
}
