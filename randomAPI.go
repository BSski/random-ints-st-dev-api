package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/BSski/RandomIntsStDevAPI/constants"
	"github.com/go-chi/chi/v5"
	"github.com/montanaflynn/stats"
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
	Result *randomAPIResult `json:"result"`
	Error  *randomAPIError  `json:"error"`
	ID     int              `json:"id"`
}

type randomAPIResult struct {
	Random randomAPIResultData `json:"random"`
}

type randomAPIResultData struct {
	Data []int `json:"data"`
}

type randomAPIError struct {
	Message string `json:"message"`
}

type partialResult struct {
	Stdev float64 `json:"stdev"`
	Data  []int   `json:"data"`
}

type randomAPIResource struct{}

func (rs randomAPIResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/mean", func(r chi.Router) {
		r.Get("/", rs.Get)
	})
	return r
}

func (rs randomAPIResource) Get(w http.ResponseWriter, r *http.Request) {
	nrOfRequestsStr := r.URL.Query().Get("requests")
	nrOfRequests, err := prepareURLParam(nrOfRequestsStr, "requests", 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	intSeqLengthStr := r.URL.Query().Get("length")
	intSeqLength, err := prepareURLParam(intSeqLengthStr, "length", 1000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	intSeqs, err := getRandomIntSeqs(r.Context(), nrOfRequests, intSeqLength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stDevsInSeqs := make([]float64, nrOfRequests)
	for i, intSeq := range intSeqs {
		roundedStDev, err := getRoundedStDev(intSeq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stDevsInSeqs[i] = roundedStDev
	}

	intSeqsSum, err := getIntSeqsSum(intSeqs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundedStDevOfSum, err := getRoundedStDev(intSeqsSum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	finalResult := make([]partialResult, nrOfRequests+1)
	for i := 0; i < nrOfRequests; i++ {
		finalResult[i] = partialResult{
			Stdev: stDevsInSeqs[i],
			Data:  intSeqs[i],
		}
	}
	finalResult[nrOfRequests] = partialResult{
		Stdev: roundedStDevOfSum,
		Data:  intSeqsSum,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(finalResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func prepareURLParam(paramStr, paramURLName string, maxVal int) (param int, err error) {
	if paramStr == "" {
		paramStr = "1"
	}

	param, err = strconv.Atoi(paramStr)
	if err != nil {
		return param, fmt.Errorf("%v param has to be an int", paramURLName)
	}

	if param <= 0 {
		return param, fmt.Errorf("%v param has to be greater than 0", paramURLName)
	}
	if param > maxVal {
		return param, fmt.Errorf("%v param has to be smaller than or equal to %v", paramURLName, maxVal)
	}
	return
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

func getRoundedStDev(intSeq []int) (roundedStDev float64, err error) {
	data := stats.LoadRawData(intSeq)
	stDev, err := stats.StandardDeviation(data)
	if err != nil {
		return
	}
	roundedStDev, err = stats.Round(stDev, 3)
	if err != nil {
		return
	}
	return
}

func getIntSeqsSum(intSeqs [][]int) (intSeqsSum []int, err error) {
	intSeqsSum = make([]int, len(intSeqs)*len(intSeqs[0]))
	for i, seq := range intSeqs {
		for j, x := range seq {
			intSeqsSum[i*len(intSeqs[0])+j] = x
		}
	}
	return
}
