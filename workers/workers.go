package workers

import (
	"fmt"
	"taskETH/request"
	"taskETH/types"
	"taskETH/utils"
	"math/big"
	"sync"
	"time"
)

func ProcessAddresses(addresses []string, blockNumberStruct *types.BlockNumber, apiKey, baseURL string) (string, *big.Int) {
	wallets := make(map[string]*big.Int)

	for _, addr := range addresses {
		startBalance := request.GetBalance(addr, blockNumberStruct, apiKey, baseURL)
		var endBalance *big.Int

		var balanceWg sync.WaitGroup
		seqChan := make(chan struct{}, 1)
		seqChan <- struct{}{}

		for i := 0; i <= 100; i++ {
			balanceWg.Add(1)
			go func(addr string, i int) {
				defer balanceWg.Done()

				<-seqChan

				pastBlock := utils.CalculatePastBlock(blockNumberStruct, i)
				balance := request.GetBalance(addr, pastBlock, apiKey, baseURL)

				if balance != nil {
					fmt.Printf("Адрес: %s, Баланс на блоке %s (с учётом %d блоков назад): %s Wei\n", addr, pastBlock.Hex, i, balance)
				}

				if i == 100 {
					endBalance = balance
				}

				time.Sleep(500 * time.Millisecond)

				seqChan <- struct{}{}
			}(addr, i)
		}
		balanceWg.Wait()
		close(seqChan)

		if startBalance != nil && endBalance != nil {
			change := utils.AbsoluteBalanceChange(startBalance, endBalance)
			wallets[addr] = change
			fmt.Printf("Адрес: %s, Баланс на начальном блоке: %s Wei, Баланс на конечном блоке: %s Wei, Изменение: %s Wei\n", addr, startBalance, endBalance, change.String())
		}
	}

	return utils.FindMaxBalanceChange(wallets)
}

func RunWorkers(addresses []string, blockNumberStruct *types.BlockNumber, numWorkers int, apiKey, baseURL string) (string, *big.Int) {
	var wg sync.WaitGroup

	addressChunks := chunkAddresses(addresses, numWorkers)
	results := make(chan struct {
		maxWallet string
		maxChange *big.Int
	}, len(addressChunks))

	for _, chunk := range addressChunks {
		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			maxWallet, maxChange := ProcessAddresses(chunk, blockNumberStruct, apiKey, baseURL)
			results <- struct {
				maxWallet string
				maxChange *big.Int
			}{maxWallet, maxChange}
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var maxWallet string
	var maxChange *big.Int

	for result := range results {
		if maxChange == nil || result.maxChange.Cmp(maxChange) > 0 {
			maxWallet = result.maxWallet
			maxChange = result.maxChange
		}
	}

	return maxWallet, maxChange
}

func chunkAddresses(addresses []string, numChunks int) [][]string {
	chunkSize := (len(addresses) + numChunks - 1) / numChunks
	chunks := make([][]string, 0, numChunks)

	for i := 0; i < len(addresses); i += chunkSize {
		end := i + chunkSize
		if end > len(addresses) {
			end = len(addresses)
		}
		chunks = append(chunks, addresses[i:end])
	}

	return chunks
}
