package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

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

type StDevOfSum struct {
	StDevOfSum float64
}

type randomAPIResource struct{}

func (rs randomAPIResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/mean", func(r chi.Router) {
		r.Get("/", rs.Get)
	})
	return r
}

func (rs randomAPIResource) requestRandomInts(ctx context.Context, intSeqLength int) (intSeq []int, err error) {
	runtime.GOMAXPROCS(1) // Random.org API guidelines prohibit simultaneous requests.
	url := "https://api.random.org/json-rpc/2/invoke"
	apiKey := os.Getenv("RANDOM_ORG_API_KEY")
	params := randomAPIParams{apiKey, intSeqLength, 1, 10}
	payload := randomAPIRequest{"2.0", "generateIntegers", params, 666}
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

func (rs randomAPIResource) validateParam(paramName string, numericParam int, maxVal int) error {
	if numericParam <= 0 {
		return fmt.Errorf("%v param has to be greater than 0", paramName)
	}

	if numericParam > maxVal {
		return fmt.Errorf("%v param has to be smaller than or equal to %v", paramName, maxVal)
	}
	return nil
}

func (rs randomAPIResource) getRandomAPIData(ctx context.Context, nrOfRequests int, intSeqLength int) (intSeqs [][]int, err error) {
	var wg sync.WaitGroup
	wg.Add(nrOfRequests)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	intSeqs = make([][]int, nrOfRequests)
	for i := 0; i < nrOfRequests; i++ {
		go func(ctx context.Context, i int) {
			defer wg.Done()

			intSeq, err := rs.requestRandomInts(ctx, intSeqLength)
			if err != nil {
				cancel()
				return
			}
			intSeqs[i] = intSeq
		}(ctx, i)
	}
	wg.Wait()
	return
}

func (rs randomAPIResource) prepareURLParam(paramStr, paramURLName string, maxVal int) (param int, err error) {
	if paramStr == "" {
		paramStr = "1"
	}
	param, err = strconv.Atoi(paramStr)
	if err != nil {
		return param, fmt.Errorf("%v param has to be greater than 0", paramURLName)
	}
	err = rs.validateParam(paramURLName, param, maxVal)
	if err != nil {
		return
	}
	return
}

func (rs randomAPIResource) getStDevOfSum(intSeqs [][]int) (roundedStDevOfSum float64, err error) {
	intSeqsSums := make([]int, len(intSeqs))
	for i, seq := range intSeqs {
		data := stats.LoadRawData(seq)
		seqSum, err := stats.Sum(data)
		if err != nil {
			return roundedStDevOfSum, err
		}
		intSeqsSums[i] = int(seqSum)
	}
	data := stats.LoadRawData(intSeqsSums)
	stDevOfSums, err := stats.StandardDeviation(data)
	if err != nil {
		return
	}
	roundedStDevOfSum, err = stats.Round(stDevOfSums, 3)
	if err != nil {
		return
	}
	return
}

func (rs randomAPIResource) getStDevs(intSeqs [][]int) (stDevsInSeqs []float64, err error) {
	stDevsInSeqs = make([]float64, len(intSeqs))
	for i, intSeq := range intSeqs {
		data := stats.LoadRawData(intSeq)
		stDev, err := stats.StandardDeviation(data)
		if err != nil {
			return stDevsInSeqs, err
		}
		roundedStDev, err := stats.Round(stDev, 3)
		if err != nil {
			return stDevsInSeqs, err
		}
		stDevsInSeqs[i] = roundedStDev
	}
	return
}

// Request Handler - GET /posts/{id} - Read a single post by :id.
func (rs randomAPIResource) Get(w http.ResponseWriter, r *http.Request) {
	nrOfRequestsStr := r.URL.Query().Get("requests")
	fmt.Println("type:", reflect.TypeOf(nrOfRequestsStr))
	nrOfRequests, err := rs.prepareURLParam(nrOfRequestsStr, "requests", 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	intSeqLengthStr := r.URL.Query().Get("length")
	intSeqLength, err := rs.prepareURLParam(intSeqLengthStr, "length", 1000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	intSeqs, err := rs.getRandomAPIData(r.Context(), nrOfRequests, intSeqLength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stDevsInSeqs, err := rs.getStDevs(intSeqs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundedStDevOfSum, err := rs.getStDevOfSum(intSeqs)
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
		Data:  intSeqs[0],
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(finalResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
