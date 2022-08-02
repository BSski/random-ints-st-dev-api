package randomintsstdev

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BSski/RandomIntsStDevAPI/randomintstdevconsts"
	"github.com/go-chi/chi/v5"
)

type partialResult struct {
	Stdev float64 `json:"stddev"`
	Data  []int   `json:"data"`
}

type RandomAPIResource struct{}

func (rs RandomAPIResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/mean", func(r chi.Router) {
		r.Get("/", rs.Get)
	})
	return r
}

func (rs RandomAPIResource) Get(w http.ResponseWriter, r *http.Request) {
	nrOfRequestsStr := r.URL.Query().Get("requests")
	nrOfRequests, err := prepareURLParam(nrOfRequestsStr, "requests", randomintstdevconsts.MAX_CONCURRENT_REQUESTS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	intSeqLengthStr := r.URL.Query().Get("length")
	intSeqLength, err := prepareURLParam(intSeqLengthStr, "length", randomintstdevconsts.MAX_SEQUENCE_LENGTH)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	intSeqsSum := getIntSeqsSum(intSeqs)

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
