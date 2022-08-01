package randomintsstdev

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/BSski/RandomIntsStDevAPI/constants"
)

type randomAPIRequest struct {
	JsonRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  randomAPIParams `json:"params"`
	ID      int             `json:"id"`
}

type randomAPIParams struct {
	APIKey string `json:"apiKey"`
	N      int    `json:"n"`
	Min    int    `json:"min"`
	Max    int    `json:"max"`
}

type randomAPIResponse struct {
	Result *randomAPIResult
	Error  *randomAPIError
	ID     int
}

type randomAPIResult struct {
	Random randomAPIResultData
}

type randomAPIResultData struct {
	Data []int
}

type randomAPIError struct {
	Message string
}

func getRandomIntSeqs(ctx context.Context, nrOfRequests int, intSeqLength int) (intSeqs [][]int, err error) {
	var wg sync.WaitGroup
	wg.Add(nrOfRequests)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	previousGOMAXPROCS := runtime.GOMAXPROCS(1) // Random.org guidelines prohibit simultaneous requests.

	intSeqs = make([][]int, nrOfRequests)
	for i := 0; i < nrOfRequests; i++ {
		go func(ctx context.Context, i int) {
			defer wg.Done()
			intSeq, err := requestRandomIntSeq(ctx, intSeqLength)
			if err != nil {
				cancel()
				return
			}
			intSeqs[i] = intSeq
		}(ctx, i)
	}
	wg.Wait()

	runtime.GOMAXPROCS(previousGOMAXPROCS)
	return
}

func requestRandomIntSeq(ctx context.Context, intSeqLength int) (intSeq []int, err error) {
	url := constants.RANDOM_API_URL
	apiKey := os.Getenv("RANDOM_ORG_API_KEY")

	params := randomAPIParams{apiKey, intSeqLength, constants.MIN_RANDOM_INT, constants.MAX_RANDOM_INT}
	payload := randomAPIRequest{"2.0", "generateIntegers", params, intSeqLength}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result randomAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return intSeq, err
	}

	if err := result.Error; err != nil {
		return intSeq, errors.New(result.Error.Message)
	}

	intSeq = result.Result.Random.Data
	return
}
