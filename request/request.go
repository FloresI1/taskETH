package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"taskETH/types"
	"math/big"
	"net/http"
	"time"
)

const (
	maxRetries = 5
	retryDelay = 500 * time.Millisecond // Базовая задержка между повторными попытками
)

func retryRequest(requestBody []byte, url string) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryDelay * time.Duration(i+1)) // Экспоненциальная задержка
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err == nil {
			return resp, nil
		}
		fmt.Printf("Ошибка при выполнении POST запроса, попытка %d: %s\n", i+1, err)
	}
	return nil, fmt.Errorf("достигнуто максимальное количество попыток (%d)", maxRetries)
}

func GetBlockNumber(apiKey, baseURL string) string {
	url := fmt.Sprintf("%s%s", baseURL, apiKey)

	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      "getblock.io",
	})
	if err != nil {
		fmt.Println("Ошибка при подготовке данных JSON:", err)
		return ""
	}

	resp, err := retryRequest(requestBody, url)
	if err != nil {
		fmt.Println("Ошибка при выполнении POST запроса:", err)
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return ""
	}

	blockNum, ok := result["result"].(string)
	if !ok {
		fmt.Println("Поле result отсутствует или имеет некорректный формат")
		return ""
	}

	return blockNum
}

func GetBlockByNumber(blockNum *types.BlockNumber, apiKey, baseURL string) (map[string]types.Transaction, error) {
	url := fmt.Sprintf("%s%s", baseURL, apiKey)

	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNum.Hex, true},
		"id":      "getblock.io",
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при подготовке данных JSON: %s", err)
	}

	resp, err := retryRequest(requestBody, url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении POST запроса: %s", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %s", err)
	}

	blockData, ok := result["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("поле result отсутствует или имеет некорректный формат")
	}

	transactionsData, ok := blockData["transactions"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("поле transactions отсутствует или имеет некорректный формат")
	}

	transactionsMap := make(map[string]types.Transaction)
	for _, tx := range transactionsData {
		txData, ok := tx.(map[string]interface{})
		if !ok {
			continue
		}

		from, fromOk := txData["from"].(string)
		to, toOk := txData["to"].(string)
		hash, hashOk := txData["hash"].(string)

		if fromOk && toOk && hashOk {
			transactionsMap[hash] = types.Transaction{
				From: from,
				To:   to,
			}
		}
	}

	return transactionsMap, nil
}

func GetBalance(address string, blockNum *types.BlockNumber, apiKey, baseURL string) *big.Int {
	url := fmt.Sprintf("%s%s", baseURL, apiKey)

	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []interface{}{address, blockNum.Hex},
		"id":      "getblock.io",
	})
	if err != nil {
		fmt.Println("Ошибка при подготовке данных JSON:", err)
		return nil
	}

	resp, err := retryRequest(requestBody, url)
	if err != nil {
		fmt.Println("Ошибка при выполнении POST запроса:", err)
		return nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return nil
	}

	resultHex, ok := result["result"].(string)
	if !ok {
		fmt.Println("Поле result отсутствует или имеет некорректный формат")
		return nil
	}

	resultDec := new(big.Int)
	if _, ok := resultDec.SetString(resultHex[2:], 16); !ok {
		fmt.Println("Ошибка при конвертации из hex в decimal")
		return nil
	}
	return resultDec
}
