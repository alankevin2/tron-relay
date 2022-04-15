package relay

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	crypto "github.com/ethereum/go-ethereum/crypto"
	address "github.com/fbsobreira/gotron-sdk/pkg/address"
	client "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/golang/protobuf/ptypes"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"gitlab.inlive7.com/crypto/tron-relay/config"
	"gitlab.inlive7.com/crypto/tron-relay/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type Relay struct {
	currentChainInfo config.ChainInfo
	supportTokens    map[string]string
	client           *client.GrpcClient
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

func (r *Relay) GetBalanceForToken(address string, symbol string) (balance *big.Int, decimal uint8, err error) {
	token := strings.ToLower(symbol)
	tokenAddress := r.supportTokens[token]
	if tokenAddress == "" {
		return nil, 0, errors.New("token not match any of supported tokens")
	}
	balance, err = r.client.TRC20ContractBalance(address, tokenAddress)
	if err != nil {
		return nil, 0, err
	}
	d, err := r.client.TRC20GetDecimals(tokenAddress)
	if err != nil {
		return nil, 0, err
	}
	decimal = uint8(d.Int64())
	return balance, decimal, nil
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

	pk, err := crypto.HexToECDSA(msg.PrivateKey)
	if err != nil {
		log.Panic("private key hex to ecdsa failed")
	}
	// TODO, discuss that should we check everytime or not?
	pub := pk.PublicKey
	pubStr := hex.EncodeToString(crypto.FromECDSAPub(&pub))
	addr := publicStrToAddress(pubStr)

	if addr != msg.Address {
		err = fmt.Errorf("address from remote does not pass local check, remote:%s, local:%s", msg.Address, addr)
		log.Panic("CreateNewAccount failed")
	}
	return msg.PrivateKey, pubStr, msg.Address, nil
}

func (r *Relay) TransferTokenUsingPrivateKey(privateKey string, data *types.TransactionRaw) (txn string, err error) {
	token := strings.ToLower(data.TokenSymbol)
	tokenAddress := r.supportTokens[token]
	if tokenAddress == "" {
		return "", errors.New("token not match any of supported tokens")
	}
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	pub := pk.PublicKey
	pubStr := hex.EncodeToString(crypto.FromECDSAPub(&pub))
	from := publicStrToAddress(pubStr)

	t, err := r.client.TRC20Send(from, data.To, tokenAddress, data.Value, int64(data.PreferredLimit))
	if err != nil {
		return "", err
	}

	rawData, err := proto.Marshal(t.Transaction.GetRawData())
	if err != nil {
		return "", err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, pk)
	if err != nil {
		return "", err
	}
	t.Transaction.Signature = append(t.Transaction.Signature, signature)
	result, err := r.client.Broadcast(t.Transaction)
	if err != nil {
		return "", err
	}
	if result.Code > 0 {
		return "", errors.New(string(result.Message))
	}
	return hex.EncodeToString(t.Txid), nil
}

func (r *Relay) QueryTransaction(txn string) (*types.TransactionState, error) {
	tInfo, err := r.client.GetTransactionInfoByID(txn)
	if err != nil {
		return nil, err
	}
	from, to, value := "", "", big.NewInt(0)
	if tInfo.ContractAddress != nil && len(tInfo.ContractAddress) > 0 {
		var isSupportToken bool
		contractAddressInHex := hex.EncodeToString(tInfo.ContractAddress)
		contractAddress := address.HexToAddress(contractAddressInHex).String()
		for _, each := range r.supportTokens {
			isSupportToken = each == contractAddress
		}
		if isSupportToken && tInfo.Log != nil && len(tInfo.Log) == 1 {
			f := bytes.Trim(tInfo.Log[0].Topics[1], "\x00")
			t := bytes.Trim(tInfo.Log[0].Topics[2], "\x00")
			from = address.HexToAddress("41" + hex.EncodeToString(f)).String()
			to = address.HexToAddress("41" + hex.EncodeToString(t)).String()
			value = new(big.Int).SetBytes(tInfo.Log[0].Data)
		}
	} else {
		t, err := r.client.GetTransactionByID(txn)
		if err != nil {
			return nil, err
		}
		if len(t.RawData.Contract) < 1 || t.RawData.Contract[0] == nil {
			return nil, errors.New("transaction rawdata contract is nil")
		}
		if t.RawData.Contract[0].Type != core.Transaction_Contract_TransferContract {
			return nil, errors.New("transaction action is out of our scope")
		}
		contract := t.RawData.Contract[0]
		var c core.TransferContract

		if err = ptypes.UnmarshalAny(contract.GetParameter(), &c); err != nil {
			return nil, err
		}
		from = address.Address(c.OwnerAddress).String()
		to = address.Address(c.ToAddress).String()
		value = big.NewInt(c.Amount)
	}

	return &types.TransactionState{
		Success:   tInfo.Result.Number() == 0,
		From:      from,
		To:        to,
		Value:     value,
		FeeLimit:  uint64(tInfo.Fee),
		Time:      uint64(tInfo.BlockTimeStamp),
		Chain:     uint16(r.currentChainInfo.ID),
		ChainName: r.currentChainInfo.Name,
	}, nil
}

func (r *Relay) GetFeeLimit() (limit uint64, err error) {
	// reference to github.com/fbsobreira/gotron-sdk in package cmd file contracts.go
	return 1000000000, nil
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
	return &Relay{client: client, currentChainInfo: c, supportTokens: p.Tokens}, nil
}

func (r *Relay) destory() {
	r.client.Conn.Close()
}

func publicStrToAddress(publicKey string) string {
	pubInByte, _ := hex.DecodeString(publicKey[2:]) // don't know why, should ask tron, reference is here:https://www.btcschools.net/tron/tron_address.php
	hashInByte := crypto.Keccak256(pubInByte)       // refernce: https://www.btcschools.net/tron/tron_address.php
	hash := hex.EncodeToString(hashInByte)          // NEVER use string(byte[]), because it will be case-sensitive, which is not expected!!!!!
	hash = "41" + hash[len(hash)-40:]               // last two byte in hex string means 40 length of digits
	return address.HexToAddress(hash).String()
}
