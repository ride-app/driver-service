//go:generate go run go.uber.org/mock/mockgen -destination ../../testing/mocks/$GOFILE -package mocks . WalletRepository

package walletrepository

import (
	"context"
	"net/http"

	walletApi "buf.build/gen/go/ride/wallet/connectrpc/go/ride/wallet/v1alpha1/walletv1alpha1connect"
	pb "buf.build/gen/go/ride/wallet/protocolbuffers/go/ride/wallet/v1alpha1"
	"connectrpc.com/connect"
	"github.com/ride-app/driver-service/pkg/config"
	"github.com/ride-app/driver-service/pkg/utils/logger"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, log logger.Logger, id string, authToken string) (*pb.Wallet, error)
}

type Impl struct {
	walletApi walletApi.WalletServiceClient
}

func New(log logger.Logger, config *config.Config) (*Impl, error) {
	log.Debug("Wallet Service Host: ", config.WalletServiceHost)
	client := walletApi.NewWalletServiceClient(
		http.DefaultClient,
		config.WalletServiceHost,
	)

	log.Info("Wallet Repository initialized")
	return &Impl{walletApi: client}, nil
}

func (r *Impl) GetWallet(ctx context.Context, log logger.Logger, id string, authToken string) (*pb.Wallet, error) {
	log.Info("Getting wallet from wallet service")
	req := connect.NewRequest(&pb.GetWalletRequest{
		Name: "users/" + id + "/wallet",
	})
	req.Header().Add("Authorization", authToken)

	res, err := r.walletApi.GetWallet(ctx, req)

	if err != nil {
		log.WithError(err).Error("Error getting wallet from wallet service")
		return nil, err
	}

	return res.Msg.Wallet, nil
}
