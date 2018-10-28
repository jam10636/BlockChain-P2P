package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type message struct {
	Ordernumber string
}

func run(block *[]Blockmember) error {
	mux := makeMuxRouter()
	s := &http.Server{
		Addr:           ":" + "8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/create", handleNewBlock).Methods("POST")
	return muxRouter
}
