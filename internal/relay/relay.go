package relay

import (
	"context"
	"log"

	client "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"gitlab.inlive7.com/crypto/tron-relay/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Relay struct {
	client *client.GrpcClient
}

var instances = make(map[config.ChainID]*Relay)

func Shared(chainID config.ChainID) *Relay {
	if instances[chainID] != nil {
		return instances[chainID]
	}

	info, err := config.RetrieveChainInfo(chainID)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	instance, err := createInstance(info)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	instances[chainID] = instance

	return instances[chainID]
}

/*
	This method is for hot-update usecase. If we manage to update the yml files,
	then destory instance to make it load the newer version of yml file.
*/
func Destory() {
	for _, v := range instances {
		v.destory()
	}
}

func (r *Relay) GetBalance(address string) (balance int64, err error) {
	acc, err := r.client.GetAccount(address)
	if err != nil {
		return 0, err
	}
	return acc.Balance, nil
}

func (r *Relay) CreateNewAccount() (privateKey string, publicKey string, publicAddress string, err error) {
	defer func() {
		if err != nil {
			privateKey = ""
			publicKey = ""
			publicAddress = ""
		}
	}()

	msg, err := r.client.Client.GenerateAddress(context.Background(), &api.EmptyMessage{}, &grpc.EmptyCallOption{})
	if err != nil {
		log.Panic("generate address failed")
	}

	return msg.PrivateKey, publicKey, msg.Address, nil
}

// ******** PRIVATE ******** //

func createInstance(c config.ChainInfo) (*Relay, error) {
	p := config.GetProviderInfo(c.ProviderFile)
	client := client.NewGrpcClient(p.URL)
	client.SetAPIKey(p.APIKey)
	// load grpc options
	opts := make([]grpc.DialOption, 0)
	// TODO here
	// if withTLS {
	// opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	// } else {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// }
	err := client.Start(opts...)
	if err != nil {
		log.Fatal(err)
	}
	return &Relay{client: client}, nil
}

func (r *Relay) destory() {
	r.client.Conn.Close()
}
