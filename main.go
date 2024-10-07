package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	"github.com/afeyisa/Cherpy/api"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}

var x = apiConfig{fileserverHits: atomic.Int32{}}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main()  {

	
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz",isHealth)	
	mux.HandleFunc("/rf",getNumberOfReq)
	mux.HandleFunc("POST /api/validate",api.ValidateChirp)
	mux.Handle("/",	x.middlewareMetricsInc(http.FileServer(http.Dir("./"))))

	
	err1 := http.ListenAndServe(":8080",mux)
	if errors.Is(err1, http.ErrServerClosed){
		fmt.Printf("Server closed\n")
	}else if err1 != nil{
		fmt.Printf("error starting server %s\n",err1)
		os.Exit(1)
	}

	
}

func getRoot( w http.ResponseWriter,r *http.Request){
	fmt.Printf("got / request\n")
	io.WriteString(w, "I am testing my website\n")

}


func isHealth (w http.ResponseWriter,r *http.Request){
	w.Header().Add("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func getNumberOfReq(w http.ResponseWriter,r *http.Request){
	w.Header().Add("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(strconv.Itoa(int(x.fileserverHits.Load()))))

}




// for Json phase
/*****************************************************************/

// json file decoder
func decodeHandler (w http.ResponseWriter,r * http.Request){
	type parameters struct {
		Name string `json:"name"`
		Age int `json:"age"`
	}

	decoder := json.NewDecoder(r.Body)
	prams := parameters{}

	err := decoder.Decode(&prams)

	if err != nil {
		log.Printf("error decoding parameteres %s",err)
		w.WriteHeader(500)
		return
	}

	// now params can be used
}


// json encoder
func encodeHandler ( w http.ResponseWriter, r *http.Request){
	type returnVals struct{
		CreatedAt time.Time `json:"created_at"`
		ID int `json:"id"`
	}
	resPBody := returnVals{
		CreatedAt:time.Now() ,
		ID: 22,
	}
	data , err := json.Marshal(resPBody)
	if err != nil {
		log.Printf("erro marshalling JSON %s",err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

/*************************************************/