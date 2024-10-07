package api

import (
	"encoding/json"
	"log"
	"net/http"
	
)

type errorResponse struct{
	Err string `json:"error"`
}

type parameters struct {
	Body string `json:"body"`
	
}

type validResponse struct{
	Valid bool `json:"valid"`
}


func ValidateChirp(w http.ResponseWriter,r *http.Request){
	decoder := json.NewDecoder(r.Body)
	prams := parameters{}
	err := decoder.Decode(&prams)

	if err != nil {
		log.Printf("error decoding parameteres %s",err)
		rbody := errorResponse{Err:"Something went Wrong"}
		data , err := json.Marshal(rbody)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(500)
		w.Write(data)
		return
		
	}

	if len(prams.Body) > 140{
		rbody := errorResponse{Err: "chirp is too long"}
		data , err := json.Marshal(rbody)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(500)
		w.Write(data)
		return

	}

	rbody := validResponse{Valid: true}
	data , err := json.Marshal(rbody)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	w.Write(data)
}