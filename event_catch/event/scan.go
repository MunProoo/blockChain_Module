package event

import (
	"context"
	"event/config"
	"fmt"
	"math/big"
	"time"

	"github.com/hacpy/go-ethereum"
	"github.com/hacpy/go-ethereum/common"
	ethTypes "github.com/hacpy/go-ethereum/core/types"
	"github.com/hacpy/go-ethereum/ethclient"
)

type Scan struct {
	config *config.Config

	FilterQuery ethereum.FilterQuery
	client      *ethclient.Client
}

func NewScan(config *config.Config, client *ethclient.Client, catchEventList []common.Hash) (*Scan, chan []ethTypes.Log, error) {
	s := &Scan{
		config: config,
		client: client,
	}

	eventLog := make(chan []ethTypes.Log, 100)
	scanCollection := common.HexToAddress("")

	go s.lookingScan(config.Node.StartBlock, scanCollection, catchEventList, eventLog)

	return s, eventLog, nil
}

func (s *Scan) lookingScan(
	startBlock int64,

	//  scan해야하는 Collection. 보통은 컨트랙트를 배포될 때, 컨트랙트의 정보를 저장하고
	// DB값에서 모듈이 실행될 때 컨트랙트의 주소를 가져오는 형식
	scanCollection common.Address,

	catchEventList []common.Hash, // 캐치해야하는 이벤트
	eventLog chan<- []ethTypes.Log,
) {
	startReadBlock, to := startBlock, uint64(0)

	s.FilterQuery = ethereum.FilterQuery{
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{catchEventList},
		FromBlock: big.NewInt(startReadBlock),
	}

	for {
		time.Sleep(time.Second * 5)

		ctx := context.Background()
		if maxBlock, err := s.client.BlockNumber(ctx); err != nil {
			fmt.Println("Get Block Number", "err", err)
			// continue
		} else {
			to = maxBlock

			if to > uint64(startReadBlock) {
				s.FilterQuery.FromBlock = big.NewInt(int64(startReadBlock))
				s.FilterQuery.ToBlock = big.NewInt(int64(to))
				tryCount := 1
				if tryCount == 3 {
					fmt.Println("Failed to get Filter", "err", err.Error())
					break
				}

			Retry:
				if logs, err := s.client.FilterLogs(ctx, s.FilterQuery); err != nil {
					// Event를 못 불러온 경우 (블록들을 1개씩 줄이고) Retry
					newTo := big.NewInt(int64(to) - 1)
					newFrom := big.NewInt(int64(startReadBlock) - 1)
					s.FilterQuery.ToBlock = newTo
					s.FilterQuery.FromBlock = newFrom

					tryCount++
					goto Retry

					// TODO -> From, to 블럭만 변형시켜서 다시 호출

				} else if len(logs) > 0 {
					eventLog <- logs
					startReadBlock = int64(to)
				}
			}
		}
	}

}
