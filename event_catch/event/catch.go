package event

import (
	"context"
	"event/config"
	"event/types"

	"github.com/hacpy/go-ethereum/common"
	ethTypes "github.com/hacpy/go-ethereum/core/types"
	"github.com/hacpy/go-ethereum/crypto"
	"github.com/hacpy/go-ethereum/ethclient"
)

type Catch struct {
	config *config.Config
	client *ethclient.Client

	needToCatchEvent map[common.Hash]types.NeedToCatchEvent //캐치 할 이벤트
}

func NewCatch(config *config.Config, client *ethclient.Client, eventChan chan []ethTypes.Log) (*Catch, error) {
	c := &Catch{
		config: config,
		client: client,
	}

	// TODO 캐치해야하는 이벤트 정의
	c.needToCatchEvent = map[common.Hash]types.NeedToCatchEvent{
		common.BytesToHash(crypto.Keccak256([]byte("Transfer(address,address,uint256)"))): {
			NeedToCatchEventFunc: c.Transfer,
		},
	}

	go c.StartToCatch(eventChan)

	return c, nil
}

func (c *Catch) Transfer(e *ethTypes.Log, tx *ethTypes.Transaction) {

}

func (c *Catch) StartToCatch(events <-chan []ethTypes.Log) {
	for event := range events {
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
					}
				}
			}

			// reverted event
			if e.Removed {
				continue
			} else if et, ok := c.needToCatchEvent[e.Topics[0]]; ok {
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

	return eventsToCatch
}
