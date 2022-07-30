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

type FinalResult struct {
	Numbers [][]int   `json:"numbers"`
	StdDevs []float64 `json:"stddevs"`
}

type randomAPIResource struct{}

func (rs randomAPIResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/mean", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", rs.Get) // GET /posts/{id} - Read a single post by :id.
	})

	return r
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "requests", r.URL.Query().Get("requests"))
		ctx = context.WithValue(ctx, "length", r.URL.Query().Get("length"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestRandomInts(intSeqLength int) (intSeq []int, err error) {
	url := "https://api.random.org/json-rpc/2/invoke"
	apiKey := os.Getenv("RANDOM_ORG_API_KEY")
	params := randomAPIParams{apiKey, intSeqLength, 1, 10}
	payload := randomAPIRequest{"2.0", "generateIntegers", params, 666}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return intSeq, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return intSeq, err
	}
	request.Header.Set("Content-Type", "application/json")
	//request.WithContext(ctx)

	httpClient := http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return intSeq, err
	}

	defer resp.Body.Close()

	var result randomAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return intSeq, err
	}

	if err := result.Error; err != nil {
		//http.Error(w, errors.New(result.Error.Message).Error(), http.StatusInternalServerError)
		return intSeq, errors.New(result.Error.Message)
	}

	intSeq = result.Result.Random.Data
	return intSeq, err
}

func (rs randomAPIResource) validateParam(numericParam int, maxVal int) error {
	if numericParam <= 0 {
		return errors.New("param has to be greater than 0")
	}

	if numericParam > maxVal {
		return fmt.Errorf("param has to be smaller than %v", maxVal)
	}

	return nil
}

// Request Handler - GET /posts/{id} - Read a single post by :id.
func (rs randomAPIResource) Get(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(1) // Random.org API guidelines prohibit simultaneous requests.

	nrOfRequestsStr := r.Context().Value("requests").(string)
	nrOfRequests, err := strconv.Atoi(nrOfRequestsStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = rs.validateParam(nrOfRequests, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	intSeqLengthStr := r.Context().Value("length").(string)
	intSeqLength, err := strconv.Atoi(intSeqLengthStr)
	if err != nil {
		return
	}
	err = rs.validateParam(intSeqLength, 1000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// http.Error(w, errors.New(result.Error.Message).Error(), http.StatusInternalServerError)

	var wg sync.WaitGroup
	wg.Add(nrOfRequests)
	intSeqs := make([][]int, nrOfRequests)
	stdDevsInSeqs := make([]float64, nrOfRequests)
	for i := 0; i < nrOfRequests; i++ {
		go func(i int) {
			defer wg.Done()

			intSeq, err := requestRandomInts(intSeqLength)
			if err != nil {
				return
			}

			intSeqs[i] = intSeq

			data := stats.LoadRawData(intSeq)
			stdDev, _ := stats.StandardDeviation(data)
			roundedStdDev, _ := stats.Round(stdDev, 3)
			stdDevsInSeqs[i] = roundedStdDev
		}(i)
	}
	wg.Wait()

	intSeqsSums := make([]int, nrOfRequests)
	for i, seq := range intSeqs {
		data := stats.LoadRawData(seq)
		seqSum, _ := stats.Sum(data)
		intSeqsSums[i] = int(seqSum)
	}
	data := stats.LoadRawData(intSeqsSums)
	stdDevOfSums, _ := stats.StandardDeviation(data)
	roundedStdDevOfSums, _ := stats.Round(stdDevOfSums, 3)

	fmt.Println("intSeqsSums:", intSeqsSums)
	fmt.Println("stdDevsInSeqs:", stdDevsInSeqs)
	fmt.Println("roundedStdDevOfSums:", roundedStdDevOfSums)

	var finalResult FinalResult
	finalResult.Numbers = intSeqs
	finalResult.StdDevs = stdDevsInSeqs

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(finalResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
