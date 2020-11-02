package payouts

import (
	"math/big"
	"os"
	"testing"

	"github.com/etclabscore/open-ethereum-pool/rpc"
	"github.com/etclabscore/open-ethereum-pool/storage"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCalculateRewards(t *testing.T) {
	blockReward, _ := new(big.Rat).SetString("5000000000000000000")
	shares := map[string]int64{"0x0": 1000000, "0x1": 20000, "0x2": 5000, "0x3": 10, "0x4": 1}
	expectedRewards := map[string]int64{"0x0": 4877996431, "0x1": 97559929, "0x2": 24389982, "0x3": 48780, "0x4": 4878}
	totalShares := int64(1025011)

	rewards := calculateRewardsForShares(shares, totalShares, blockReward)
	expectedTotalAmount := int64(5000000000)

	totalAmount := int64(0)
	for login, amount := range rewards {
		totalAmount += amount

		if expectedRewards[login] != amount {
			t.Errorf("Amount for %v must be equal to %v vs %v", login, expectedRewards[login], amount)
		}
	}
	if totalAmount != expectedTotalAmount {
		t.Errorf("Total reward must be equal to block reward in Shannon: %v vs %v", expectedTotalAmount, totalAmount)
	}
}

func TestChargeFee(t *testing.T) {
	orig, _ := new(big.Rat).SetString("5000000000000000000")
	value, _ := new(big.Rat).SetString("5000000000000000000")
	expectedNewValue, _ := new(big.Rat).SetString("3750000000000000000")
	expectedFee, _ := new(big.Rat).SetString("1250000000000000000")
	newValue, fee := chargeFee(orig, 25.0)

	if orig.Cmp(value) != 0 {
		t.Error("Must not change original value")
	}
	if newValue.Cmp(expectedNewValue) != 0 {
		t.Error("Must charge and deduct correct fee")
	}
	if fee.Cmp(expectedFee) != 0 {
		t.Error("Must charge fee")
	}
}

func TestWeiToShannonInt64(t *testing.T) {
	wei, _ := new(big.Rat).SetString("1000000000000000000")
	origWei, _ := new(big.Rat).SetString("1000000000000000000")
	shannon := int64(1000000000)

	if weiToShannonInt64(wei) != shannon {
		t.Error("Must convert to Shannon")
	}
	if wei.Cmp(origWei) != 0 {
		t.Error("Must charge original value")
	}
}

func TestGetUncleReward(t *testing.T) {
	rewards := make(map[int64]string)
	expectedRewards := map[int64]string{
		1: "4000000000000000000",
		2: "0", //previous blocks not rewarded
		3: "0",
		4: "0",
		5: "0",
		6: "0",
	}
	for i := int64(1); i < 7; i++ {
		rewards[i] = getUncleReward(1, i+1).String()
	}
	for i, reward := range rewards {
		if expectedRewards[i] != rewards[i] {
			t.Errorf("Incorrect uncle reward for %v, expected %v vs %v", i, expectedRewards[i], reward)
		}
	}

	// Year 1
	rewardsYear1 := make(map[int64]string)
	expectedRewardsYear1 := map[int64]string{
		358363: "3500000000000000000",
		358364: "0", //previous blocks not rewarded
		358365: "0",
		358366: "0",
		358367: "0",
		358368: "0",
	}
	for i := int64(358363); i < 358363+6; i++ {
		rewardsYear1[i] = getUncleReward(358363, i+1).String()
	}
	for i, reward := range rewardsYear1 {
		if expectedRewardsYear1[i] != rewardsYear1[i] {
			t.Errorf("Incorrect uncle reward for %v, expected %v vs %v", i, expectedRewardsYear1[i], reward)
		}
	}

	// Year 2
	expectedRewardsYear2 := "3000000000000000000"
	rewardsYear2 := getUncleReward(716727, 716727+1).String()
	if expectedRewardsYear2 != rewardsYear2 {
		t.Errorf("Incorrect uncle reward, expected %v vs %v", expectedRewardsYear2, rewardsYear2)
	}

	// Year 3
	expectedRewardsYear3 := "2500000000000000000"
	rewardsYear3 := getUncleReward(1075090, 1075090+1).String()
	if expectedRewardsYear3 != rewardsYear3 {
		t.Errorf("Incorrect uncle reward, expected %v vs %v", expectedRewardsYear3, rewardsYear3)
	}

	// Year 4
	expectedRewardsYear4 := "2000000000000000000"
	rewardsYear4 := getUncleReward(1433454, 1433454+1).String()
	if expectedRewardsYear4 != rewardsYear4 {
		t.Errorf("Incorrect uncle reward, expected %v vs %v", expectedRewardsYear4, rewardsYear4)
	}

	// Year 5
	expectedRewardsYear5 := "1500000000000000000"
	rewardsYear5 := getUncleReward(1791818, 1791818+1).String()
	if expectedRewardsYear5 != rewardsYear5 {
		t.Errorf("Incorrect uncle reward, expected %v vs %v", expectedRewardsYear5, rewardsYear5)
	}

	// Year 6
	expectedRewardsYear6 := "1000000000000000000"
	rewardsYear6 := getUncleReward(2150181, 2150181+1).String()
	if expectedRewardsYear6 != rewardsYear6 {
		t.Errorf("Incorrect uncle reward, expected %v vs %v", expectedRewardsYear6, rewardsYear6)
	}

	// Year 7
	expectedRewardsYear7 := "500000000000000000"
	rewardsYear7 := getUncleReward(2508545, 2508545+1).String()
	if expectedRewardsYear7 != rewardsYear7 {
		t.Errorf("Incorrect uncle reward, expected %v vs %v", expectedRewardsYear7, rewardsYear7)
	}
}

func TestMatchCandidate(t *testing.T) {
	gethBlock := &rpc.GetBlockReply{Hash: "0x12345A", Nonce: "0x1A"}
	parityBlock := &rpc.GetBlockReply{Hash: "0x12345A", SealFields: []string{"0x0A", "0x1A"}}
	candidate := &storage.BlockData{Nonce: "0x1a"}
	orphan := &storage.BlockData{Nonce: "0x1abc"}

	if !matchCandidate(gethBlock, candidate) {
		t.Error("Must match with nonce")
	}
	if !matchCandidate(parityBlock, candidate) {
		t.Error("Must match with seal fields")
	}
	if matchCandidate(gethBlock, orphan) {
		t.Error("Must not match with orphan with nonce")
	}
	if matchCandidate(parityBlock, orphan) {
		t.Error("Must not match orphan with seal fields")
	}

	block := &rpc.GetBlockReply{Hash: "0x12345A"}
	immature := &storage.BlockData{Hash: "0x12345a", Nonce: "0x0"}
	if !matchCandidate(block, immature) {
		t.Error("Must match with hash")
	}
}
