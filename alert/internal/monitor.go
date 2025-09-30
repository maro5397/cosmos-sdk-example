package internal

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	cometbft "github.com/cometbft/cometbft/rpc/client/http"

	"alert/pkg"
)

func checkSyncStatus(
	ctx context.Context,
	client *cometbft.HTTP,
	node pkg.Node,
	config *pkg.Config,
	notifier Notifier,
	lastHeight int64,
	lastChange time.Time,
) (int64, time.Time) {
	statusContext, cancel := context.WithTimeout(ctx, config.RPCTimeout)
	defer cancel()
	status, err := client.Status(statusContext)
	if err != nil {
		log.Printf("[%s] Status 에러: %v", node.Name, err)
		return lastHeight, lastChange
	}
	h := status.SyncInfo.LatestBlockHeight
	catching := status.SyncInfo.CatchingUp
	if catching {
		if lastHeight == h {
			if lastChange.IsZero() {
				lastChange = time.Now()
			}
			if time.Since(lastChange) >= config.StopDetectWindow {
				networkContext, cancel := context.WithTimeout(ctx, config.RPCTimeout)
				defer cancel()
				networkInformation, err := client.NetInfo(networkContext)
				if err == nil && len(networkInformation.Peers) == 0 {
					message := fmt.Sprintf("[%s] 동기화 정체 감지 (height=%d, peers=0)", node.Name, h)
					_ = notifier.Notify(ctx, message)
					log.Println(message)
				}
			}
		} else {
			lastHeight, lastChange = h, time.Now()
		}
	} else {
		lastHeight, lastChange = h, time.Now()
	}
	return lastHeight, lastChange
}

func checkMissedSignature(
	ctx context.Context,
	client *cometbft.HTTP,
	node pkg.Node,
	config *pkg.Config,
	notifier Notifier,
	address []byte,
) {
	blockContext, cancel := context.WithTimeout(ctx, config.RPCTimeout)
	defer cancel()
	block, err := client.Block(blockContext, nil)
	if err != nil || block == nil || block.Block == nil {
		if err != nil {
			log.Printf("[%s] Block 정보 조회 에러: %v", node.Name, err)
		}
		return
	}
	found := false
	for _, signature := range block.Block.LastCommit.Signatures {
		if signature.BlockIDFlag == 2 && bytes.Equal(signature.ValidatorAddress.Bytes(), address) {
			found = true
			break
		}
	}
	if !found {
		message := fmt.Sprintf("[%s] 블록 %d 서명 누락 감지 (LastCommit은 %d에 대한 투표)",
			node.Name, block.Block.Height, block.Block.Height-1)
		_ = notifier.Notify(ctx, message)
		log.Println(message)
	}
}

func MonitorNode(ctx context.Context, node pkg.Node, config *pkg.Config, notifier Notifier) {
	client, err := cometbft.New(node.RPC)
	if err != nil {
		log.Printf("[%s] RPC 연결 실패: %v", node.Name, err)
		return
	}

	address, err := hex.DecodeString(node.ValidatorAddress)
	if err != nil {
		log.Printf("[%s] 검증자 주소(HEX) 파싱 실패: %v", node.Name, err)
		return
	}

	var lastHeight int64
	var lastChange time.Time
	tick := time.NewTicker(config.PollInterval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			lastHeight, lastChange = checkSyncStatus(ctx, client, node, config, notifier, lastHeight, lastChange)
			checkMissedSignature(ctx, client, node, config, notifier, address)
		}
	}
}
