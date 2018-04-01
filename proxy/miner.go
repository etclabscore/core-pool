package proxy

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/etclabscore/go-etchash"
	"github.com/ethereum/go-ethereum/common"

	"github.com/etclabscore/core-pool/util"
)

var (
	maxUint256                             = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
	ecip1099FBlockClassic uint64           = 11700000 // classic mainnet
	ecip1099FBlockMordor  uint64           = 2520000  // mordor
	uip1FEpoch            uint64           = 22       // ubiq mainnet
	hasher                *etchash.Etchash = nil
)

func (s *ProxyServer) processShare(login, id, ip string, t *BlockTemplate, params []string, stratum bool) (bool, bool, error) {
	if hasher == nil {
		if s.config.Network == "classic" {
			hasher = etchash.New(&ecip1099FBlockClassic, nil)
		} else if s.config.Network == "mordor" {
			hasher = etchash.New(&ecip1099FBlockMordor, nil)
		} else if s.config.Network == "ubiq" {
			hasher = etchash.New(nil, &uip1FEpoch)
		} else if s.config.Network == "ethereum" || s.config.Network == "ropsten" {
			hasher = etchash.New(nil, nil)
		} else {
			// unknown network
			log.Printf("Unknown network configuration %s", s.config.Network)
			return false, false, nil
		}
	}
	nonceHex := params[0]
	hashNoNonce := params[1]
	mixDigest := params[2]
	nonce, _ := strconv.ParseUint(strings.Replace(nonceHex, "0x", "", -1), 16, 64)
	shareDiff := s.config.Proxy.Difficulty

	var result common.Hash
	if stratum {
		hashNoNonceTmp := common.HexToHash(params[2])

		mixDigestTmp, hashTmp := hasher.Compute(t.Height, hashNoNonceTmp, nonce)
		params[1] = hashNoNonceTmp.Hex()
		params[2] = mixDigestTmp.Hex()
		hashNoNonce = params[1]
		result = hashTmp
	} else {
		hashNoNonceTmp := common.HexToHash(hashNoNonce)
		mixDigestTmp, hashTmp := hasher.Compute(t.Height, hashNoNonceTmp, nonce)

		// check mixDigest
		if mixDigestTmp.Hex() != mixDigest {
			return false, false, nil
		}
		result = hashTmp
	}

	// Block "difficulty" is BigInt
	// NiceHash "difficulty" is float64 ...
	// diffFloat => target; then: diffInt = 2^256 / target
	shareDiffCalc := util.TargetHexToDiff(result.Hex()).Int64()
	shareDiffFloat := util.DiffIntToFloat(shareDiffCalc)
	if shareDiffFloat < 0.0001 {
		log.Printf("share difficulty too low, %f < %d, from %v@%v", shareDiffFloat, t.Difficulty, login, ip)
		return false, false, nil
	}

	h, ok := t.headers[hashNoNonce]
	if !ok {
		log.Printf("Stale share from %v@%v", login, ip)
		return false, false, nil
	}

	if s.config.Proxy.Debug {
		log.Printf("Difficulty pool/block/share = %d / %d / %d(%f) from %v@%v", shareDiff, t.Difficulty, shareDiffCalc, shareDiffFloat, login, ip)
	}

	// check share difficulty
	shareTarget := new(big.Int).Div(maxUint256, big.NewInt(shareDiff))
	if result.Big().Cmp(shareTarget) > 0 {
		return false, false, nil
	}

	// check target difficulty
	target := new(big.Int).Div(maxUint256, big.NewInt(h.diff.Int64()))
	if result.Big().Cmp(target) <= 0 {
		ok, err := s.rpc().SubmitBlock(params)
		if err != nil {
			log.Printf("Block submission failure at height %v for %v: %v", h.height, t.Header, err)
		} else if !ok {
			log.Printf("Block rejected at height %v for %v", h.height, t.Header)
			return false, false, nil
		} else {
			s.fetchBlockTemplate()
			exist, err := s.backend.WriteBlock(login, id, params, shareDiff, h.diff.Int64(), h.height, s.hashrateExpiration)
			if exist {
				return true, false, nil
			}
			if err != nil {
				log.Println("Failed to insert block candidate into backend:", err)
			} else {
				log.Printf("Inserted block %v to backend", h.height)
			}
			log.Printf("Block found by miner %v@%v at height %d", login, ip, h.height)
		}
	} else {
		// check hashrate limit
		if s.config.Proxy.HashLimit > 0 {
			currentHashrate, _ := s.backend.GetCurrentHashrate(login)

			if s.config.Proxy.HashLimit > 0 && currentHashrate > s.config.Proxy.HashLimit {
				err := fmt.Errorf("hashLimit exceed: %v(current) > %v(hashLimit)", currentHashrate, s.config.Proxy.HashLimit)
				log.Println("Failed to insert share data into backend:", err)
				return false, false, err
			}
		}

		exist, err := s.backend.WriteShare(login, id, params, shareDiff, h.height, s.hashrateExpiration)
		if exist {
			return true, false, nil
		}
		if err != nil {
			log.Println("Failed to insert share data into backend:", err)
		}
	}
	return false, true, nil
}
