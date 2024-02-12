package event

import (
	"context"
	"event/config"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

	//contract 배포 시 contract의 address -> DB에 저장하거나 config를 통해서 제어하면 좋겠네
	scanCollection := common.HexToAddress("0xd721d1E5Df6cf45AB88F3F834c08a361390898F7")

	go s.lookingScan(config.Node.StartBlock, scanCollection, catchEventList, eventLog)

	return s, eventLog, nil
}

// Contract : 0xd721d1E5Df6cf45AB88F3F834c08a361390898F7
// 배포 블록 : 45833213
// Mint 블록 : 45834926
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
		Addresses: []common.Address{scanCollection},
		Topics:    [][]common.Hash{catchEventList},
		FromBlock: big.NewInt(startReadBlock),
	}

	fmt.Println("FromBlcok : ", s.FilterQuery.FromBlock)

	for {
		ctx := context.Background()
		time.Sleep(time.Millisecond * 100)

		if maxBlock, err := s.client.BlockNumber(ctx); err != nil {
			fmt.Println("Get Block Number", "err", err.Error())
			// Get Block Number err The method platon_blockNumber does not exist/is not available
			// 해당 에러인 경우 remix에서 컨트랙트 배포를 infura를 통해 만든 API서버로 배포하지 않아서 임을 의심하도록.

			// ethereum 메서드를 사용해야하는데, ethereum classic 메서드를 사용해서 라고 한다..? 뭐지

			/* 의존성 문제였다.
			"github.com/hacpy/go-ethereum" 이 패키지를 의존하고 있었는데
			"github.com/ethereum/go-ethereum" 이 패키지를 의존해야함.
			*/

		} else {
			to = maxBlock

			if to > uint64(startReadBlock) {

				fmt.Println("from block", s.FilterQuery.FromBlock, "to block", to)

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
