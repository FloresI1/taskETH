package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"taskETH/request"
	"taskETH/utils"
	"taskETH/workers"
	"log"
	"os"
	"strconv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
}

func main() {
	loadEnv()

	fmt.Println("Start")

	apiKey := os.Getenv("GETBLOCK_API_KEY")
	baseURL := os.Getenv("GETBLOCK_BASE_URL")

	blockNum := request.GetBlockNumber(apiKey, baseURL)
	if blockNum == "" {
		fmt.Println("Ошибка: не удалось получить номер блока")
		return
	}

	blockNumberStruct := utils.ConvertBlockNumber(blockNum)

	transactions, err := request.GetBlockByNumber(blockNumberStruct, apiKey, baseURL)
	if err != nil {
		fmt.Println("Ошибка при получении данных о блоке:", err)
		return
	}

	addresses := utils.ExtractAddresses(transactions)

	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		fmt.Printf("Ошибка конвертации числа работников: %v\n", err)
		return
	}

	maxWallet, maxChange := workers.RunWorkers(addresses, blockNumberStruct, numWorkers, apiKey, baseURL)

	fmt.Printf("Кошелек с максимальным изменением: %s, изменение: %s Wei\n", maxWallet, maxChange.String())
}
