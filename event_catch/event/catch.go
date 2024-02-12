package event

import (
	"context"
	"event/config"
	"event/repository"
	"event/types"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Catch struct {
	config *config.Config

	client     *ethclient.Client
	repository *repository.Repository

	needToCatchEvent map[common.Hash]types.NeedToCatchEvent //캐치 할 이벤트
}

func NewCatch(config *config.Config, client *ethclient.Client, repository *repository.Repository) (*Catch, error) {
	c := &Catch{
		config:     config,
		client:     client,
		repository: repository,
	}

	// TODO 캐치해야하는 이벤트 리스트 정의
	c.needToCatchEvent = map[common.Hash]types.NeedToCatchEvent{
		common.BytesToHash(crypto.Keccak256([]byte("Transfer(address,address,uint256)"))): {
			NeedToCatchEventFunc: c.Transfer,
		},
	}

	// go c.StartToCatch(eventChan)

	return c, nil
}

func (c *Catch) Transfer(e *ethTypes.Log, tx *ethTypes.Transaction) {
	fmt.Println("Transfer 캐치했습니다")

	// 인덱싱이 안되어있으면
	// 1번쨰 필드 e.Data[:0x20]
	// 2번째 필드 e.Data[0x20:0x40]
	// 인덱싱이 되어있으므로 Topic 이렇게 사용 가능
	from := common.BytesToAddress(e.Topics[1][:])
	to := common.BytesToAddress(e.Topics[2][:])
	tokenID := new(big.Int).SetBytes(e.Topics[3][:])

	chainID, _ := c.client.ChainID(context.Background())
	sender, _ := ethTypes.Sender(ethTypes.NewLondonSigner(chainID), tx)

	var err error
	if err = c.repository.UpsertTxEvent(from, to, sender, tokenID, e.TxHash.Hex()); err != nil {
		fmt.Println("Failed to upsert tx event", "err", err)
	}
	if err = c.repository.UpsertNFTEvent(tokenID, to); err != nil {
		fmt.Println("Failed to upsert nft event", "err", err)
	}

}

func (c *Catch) StartToCatch(events <-chan []ethTypes.Log) {
	for event := range events {
		fmt.Println("캐치중입니다")
		ctx := context.Background()

		txList := make(map[common.Hash]*ethTypes.Transaction)

		for _, e := range event {

			hash := e.TxHash

			// 캐치한 이벤트 중 원하는 이벤트 해시가 있으면
			if _, ok := txList[hash]; !ok {
				if tx, pending, err := c.client.TransactionByHash(ctx, hash); err == nil {
					// 블록에 트랜잭션이 들어간 상태만 처리
					if !pending {
						txList[hash] = tx
					} else {
						fmt.Println("no pending transaction")
					}
				} else {
					fmt.Println(err)
				}
			}

			// reverted event
			// e.Topic -> transaction Hash

			if e.Removed {
				continue
			} else if et, ok := c.needToCatchEvent[e.Topics[0]]; !ok {
				fmt.Println("reverted")
				// 내가 캐치하고픈 이벤트가 아닌경우 TODO Log
			} else {
				et.NeedToCatchEventFunc(&e, txList[hash])
			}
		}
	}

}

func (c *Catch) GetEventToCatch() []common.Hash {
	eventsToCatch := make([]common.Hash, 0)
	for e := range c.needToCatchEvent {
		eventsToCatch = append(eventsToCatch, e)
	}

	fmt.Println("get Events to Catch", eventsToCatch)

	return eventsToCatch
}
