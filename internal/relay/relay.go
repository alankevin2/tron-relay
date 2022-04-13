package relay

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	address "github.com/fbsobreira/gotron-sdk/pkg/address"
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

	// god knows why wallet client api do not return Public key....
	msg, err := r.client.Client.GenerateAddress(context.Background(), &api.EmptyMessage{}, &grpc.EmptyCallOption{})
	if err != nil {
		log.Panic("generate address failed")
	}

	// TODO, discuss that should we check everytime or not?
	pub := public(msg.PrivateKey)             // thank you, Stackoverflow
	pubInByte, _ := hex.DecodeString(pub[2:]) // don't know why, should ask tron, reference is here:https://www.btcschools.net/tron/tron_address.php
	hashInByte := crypto.Keccak256(pubInByte) // refernce: https://www.btcschools.net/tron/tron_address.php
	hash := hex.EncodeToString(hashInByte)    // NEVER use string(byte[]), because it will be case-sensitive, which is not expected!!!!!
	hash = "41" + hash[len(hash)-40:]         // last two byte in hex string means 40 length of digits
	addr := address.HexToAddress(hash).String()

	if addr != msg.Address {
		err = fmt.Errorf("address from remote does not pass local check, remote:%s, local:%s", msg.Address, addr)
		log.Panic("CreateNewAccount faileds")
	}
	return msg.PrivateKey, pub, msg.Address, nil
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

func public(privateKey string) (publicKey string) {
	var e ecdsa.PrivateKey
	e.D, _ = new(big.Int).SetString(privateKey, 16)
	e.PublicKey.Curve = secp256k1.S256()
	e.PublicKey.X, e.PublicKey.Y = e.PublicKey.Curve.ScalarBaseMult(e.D.Bytes())
	return fmt.Sprintf("%x", elliptic.Marshal(secp256k1.S256(), e.X, e.Y))
}
