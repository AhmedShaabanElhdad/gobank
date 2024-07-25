package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
}

func NewApiServer(listentAddr string) *APIServer {
	return &APIServer{
		listenAddress: listentAddr,
	}
}

// this is use to start our server
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleSpecificAccount))

	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleSpecificAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id := mux.Vars(r)["id"]
		i, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("id %s not found ", i)
		}
		writeJson(w, http.StatusOK, &Account{
			ID: i,
		})
		return s.handleGetAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Ahmed", "Shaban")
	return writeJson(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// router.HandleFunc("/account", s.handleAccount) --> s.handleAccount return error not in func(http.ResponseWriter, *http.Request) without any return
// we need to convert our function to http.HandlerFunc
// firest we need to make type with function signature to use it in decorator function makeHttpHandlerFunc
// we keep handleAccount(w http.ResponseWriter, r *http.Request) with reutrn error despite remove error will solve the problem without need to make handle function
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// here we need remove the error to be HandlerFunc
// this is the place where we will handle error
func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error to json also
			writeJson(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}

	}
}

// Now as we make restful api will make json for response
// Encode function return error
func writeJson(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}
