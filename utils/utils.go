package utils

import (
	"math/big"
	"strings"

	"taskETH/types"
)

func ConvertBlockNumber(hexBlockNum string) *types.BlockNumber {
	intBlockNum := new(big.Int)
	intBlockNum.SetString(strings.TrimPrefix(hexBlockNum, "0x"), 16)

	return &types.BlockNumber{
		Hex: hexBlockNum,
		Int: intBlockNum,
	}
}

func ExtractAddresses(transactions map[string]types.Transaction) []string {
	addressSet := make(map[string]struct{})
	for _, tx := range transactions {
		addressSet[tx.From] = struct{}{}
		addressSet[tx.To] = struct{}{}
	}

	addresses := make([]string, 0, len(addressSet))
	for addr := range addressSet {
		addresses = append(addresses, addr)
	}
	return addresses
}

func CalculatePastBlock(currentBlock *types.BlockNumber, blocksAgo int) *types.BlockNumber {
	pastBlockInt := new(big.Int).Sub(currentBlock.Int, big.NewInt(int64(blocksAgo)))
	pastBlockHex := "0x" + pastBlockInt.Text(16)
	return &types.BlockNumber{
		Hex: pastBlockHex,
		Int: pastBlockInt,
	}
}

func AbsoluteBalanceChange(startBalance, endBalance *big.Int) *big.Int {
	change := new(big.Int)
	change.Abs(change.Sub(endBalance, startBalance))
	return change
}

func FindMaxBalanceChange(wallets map[string]*big.Int) (string, *big.Int) {
	var maxWallet string
	var maxChange *big.Int

	for wallet, change := range wallets {
		if maxChange == nil || change.Cmp(maxChange) > 0 {
			maxWallet = wallet
			maxChange = change
		}
	}

	return maxWallet, maxChange
}
