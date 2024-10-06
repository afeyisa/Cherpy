package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)




func main()  {

	mux := http.NewServeMux()
	mux.Handle("/",http.FileServer(http.Dir("./")))
	mux.HandleFunc("/healthz",isHealth)	
	
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