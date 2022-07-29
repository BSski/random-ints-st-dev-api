package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type FinalResult struct {
	StdDev  []int `json:"stddev"`
	Numbers []int `json:"numbers"`
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

// Request Handler - GET /posts/{id} - Read a single post by :id.
func (rs randomAPIResource) Get(w http.ResponseWriter, r *http.Request) {
	//requests := r.Context().Value("requests").(string)
	length := r.Context().Value("length").(string)

	link := fmt.Sprintf("https://www.random.org/integers/?num=%v&min=1&max=6&col=1&base=10&format=plain&rnd=new", length)
	resp, err := http.Get(link)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	bodyStr := strings.Fields(bodyString)

	var nrs = []int{}

	for _, i := range bodyStr {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		nrs = append(nrs, j)
	}
	fmt.Println(nrs)

	fmt.Println(reflect.TypeOf(nrs[0]))

	var result FinalResult
	result.Numbers = nrs
	result.StdDev = []int{33, 4, 6, 7}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
