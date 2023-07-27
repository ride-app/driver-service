//go:generate go run go.uber.org/mock/mockgen -destination ../../mocks/$GOFILE -package mocks . WalletRepository

package walletrepository

import (
	"context"
	"net/http"

	walletApi "buf.build/gen/go/ride/wallet/bufbuild/connect-go/ride/wallet/v1alpha1/walletv1alpha1connect"
	pb "buf.build/gen/go/ride/wallet/protocolbuffers/go/ride/wallet/v1alpha1"
	"github.com/bufbuild/connect-go"
	"github.com/ride-app/driver-service/config"
	"github.com/sirupsen/logrus"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, id string, authToken string) (*pb.Wallet, error)
}

type Impl struct {
	walletApi walletApi.WalletServiceClient
}

func New() (*Impl, error) {
	logrus.Debug("Wallet Service Host: ", config.Env.Wallet_Service_Host)
	client := walletApi.NewWalletServiceClient(
		http.DefaultClient,
		config.Env.Wallet_Service_Host,
	)

	logrus.Info("WalletRepository initialized")
	return &Impl{walletApi: client}, nil
}

func (r *Impl) GetWallet(ctx context.Context, id string, authToken string) (*pb.Wallet, error) {
	logrus.Info("Getting wallet from wallet service")
	req := connect.NewRequest(&pb.GetWalletRequest{
		Name: "users/" + id + "/wallet",
	})
	req.Header().Add("Authorization", authToken)

	res, err := r.walletApi.GetWallet(ctx, req)

	if err != nil {
		logrus.Error("Error getting wallet from wallet service: ", err)
		return nil, err
	}

	return res.Msg.Wallet, nil
}
