package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

////////////////////////////////// Server /////////////////////////////////////

// server should has listen address will use in mux as usual
// also give you the power to change with dev or other environment
// also server will need db store to make CRUD
type APIServer struct {
	listenAddress string
	store         Storage
}

func NewApiServer(listentAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listentAddr,
		store:         store,
	}
}

////////////////////////////////// Router /////////////////////////////////////

// this is use to start our server
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleSpecificAccount))

	router.HandleFunc("/accounts", makeHttpHandlerFunc(s.handleGetAccounts))

	// here is not a good approach to show accountNumber related to security and privacy and can be show in browser history
	// router.HandleFunc("/transfer/{accountNumber}", makeHttpHandlerFunc(s.handleTransfer))

	router.HandleFunc("/transfer", makeHttpHandlerFunc(s.handleTransfer))

	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "PUT" {
		return s.handleUpdateccount(w, r)
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
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, account)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		sqlRes, err := s.store.GetAccounts()
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, sqlRes)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(CreateAccountRequest) // -->  return pointer and will use it in json and other stuff so better to use pointer
	// createAccountRequest := CreateAccountRequest{}    --> return instance not pointer so will need to use &createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}
	defer r.Body.Close()
	account := NewAccount(
		createAccountRequest.FirstName,
		createAccountRequest.LastName,
	)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return writeJson(w, http.StatusCreated, account)
}

func (s *APIServer) handleUpdateccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		transferRequest := new(TransferRequest)
		if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
			return err
		}
		defer r.Body.Close()

		return writeJson(w, http.StatusCreated, transferRequest)
	}
	return fmt.Errorf("method not allowed %s", r.Method)

}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func permissionDenided(w http.ResponseWriter) {
	writeJson(w, http.StatusForbidden, ApiError{Error: "Permission Denied"})
}

////////////////////////////////// Token Helper /////////////////////////////////////

// this is another decorator start from check jwt -> handlerFunc
func withJwtAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// imagine we will add in header in x-jwt-token
		tokenString := r.Header.Get("x-jwt-token")
		token, err := verifyToken(tokenString)
		if err != nil {
			permissionDenided(w)
			return
		}
		if !token.Valid {
			permissionDenided(w)
			return
		}

		// this is one of tricky way to get propery as mapClaims
		AccountId, err := getID(r)
		if err != nil {
			writeJson(w, http.StatusBadRequest, "Missing required Data")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		// this how you can get the type from claims
		// be carful if u use your own claims as at the end will be map so becarful of the type
		// panic(reflect.TypeOf(claims["AccountId"]))  // return float64
		if int(claims["AccountId"].(float64)) != AccountId {
			permissionDenided(w)
			return
		}

		handlerFunc(w, r)
	}
}

// normal will take user to build claims but we will treat account as user for now
func createJwtAuth(account *Account) (string, error) {
	// Create the Claims
	//  use &jwt.RegisteredClaims is create Claim with Standard value
	claims := &jwt.MapClaims{
		"ExpiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"AccountId": account.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	mySigningKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(mySigningKey))
}

func verifyToken(tokenString string) (*jwt.Token, error) {

	// need to move this outside the environment to for example github secret
	// we can use export JWT_SECRET=bank5454 for now
	// var secretKey = []byte("secret-key")
	secrets := os.Getenv("JWT_SECRET")

	// we need to convert the string to Token object using Parse which will return (*jwt.Token, error)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secrets), nil
	})
}

////////////////////////////////// Middleware  /////////////////////////////////////

// router.HandleFunc("/account", s.handleAccount) --> s.handleAccount return error not in func(http.ResponseWriter, *http.Request) without any return
// we need to convert our function to http.HandlerFunc
// firest we need to make type with function signature to use it in decorator function makeHttpHandlerFunc
// we keep handleAccount(w http.ResponseWriter, r *http.Request) with reutrn error despite remove error will solve the problem without need to make handle function
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
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
// will use it in correct response or error
func writeJson(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}
