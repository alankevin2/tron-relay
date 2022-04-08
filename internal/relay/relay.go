package relay

import (
	"log"

	client "github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Relay struct {
	client *client.GrpcClient
}

var instances = make(map[uint8]*Relay)

func Shared(chainID uint8) *Relay {
	instance := instances[chainID]
	if instance == nil {
		var apiKey string
		var endPoint string
		if chainID == 1 {
			apiKey = "f1e478d5-f502-4121-8c1e-0b8ac3f47d8b"
			endPoint = "grpc.trongrid.io:50051"
		} else {
			apiKey = "f1e478d5-f502-4121-8c1e-0b8ac3f47d8b"
			endPoint = "grpc.shasta.trongrid.io:50051"
		}
		c := client.NewGrpcClient(endPoint)
		c.SetAPIKey(apiKey)
		// load grpc options
		opts := make([]grpc.DialOption, 0)
		// TODO here
		// if withTLS {
		// opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
		// } else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// }
		err := c.Start(opts...)
		if err != nil {
			log.Fatal(err)
		}
		instance = &Relay{client: c}
		instances[chainID] = instance
	}
	return instance
}

func (r *Relay) GetBalance(address string) (balance int64, err error) {
	acc, err := r.client.GetAccount(address)
	if err != nil {
		return 0, err
	}
	return acc.Balance, nil
}
