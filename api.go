package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))
	http.Handle("/", router)

	log.Println("JSON API running on port: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method is not allowed %s", r.Method)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		invalidMethod(w)
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return err
	}

	account, err := s.store.GetAccountByEmail(loginReq.Email)
	if err != nil {
		permissionDenied(w)
	}

	if !account.validatePassword(loginReq.Password) {
		return fmt.Errorf("Incorrect Email or Password")
	}

	token, err := createJWT(account)
	if err != nil {
		return err
	}

	response := LoginResponse{
		AccountNumber: account.AccountNumber,
		Token:         token,
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {

		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// GET /account
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account, err := NewAccount(createAccountReq.Email, createAccountReq.Password, createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.PhoneNumber)
	if err != nil {
		return err
	}

	if err := validateAccountInfo(account); err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT Token: ", tokenString)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferReq)
}

// JWT Functions
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":     15000,
		"accountNumber": account.AccountNumber,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func withJWTAuth(handlerFunc http.HandlerFunc, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			log.Println("Invalid Token")
			return
		}

		if !token.Valid {
			permissionDenied(w)
			log.Println("Invalid Token")
			return
		}

		userId, err := getID(r)
		if err != nil {
			permissionDenied(w)
			log.Println("No id in params")
			return
		}
		account, err := store.GetAccountByID(userId)
		if err != nil {
			permissionDenied(w)
			log.Println("Invalid account number")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.AccountNumber != int64(claims["accountNumber"].(float64)) {
			log.Println("Account numbers are not equal")
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

// validate functions
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func validateAccountInfo(info *Account) error {
	if strings.TrimSpace(info.Email) == "" || strings.TrimSpace(info.FirstName) == "" || strings.TrimSpace(info.LastName) == "" || strings.TrimSpace(info.PhoneNumber) == "" {
		return fmt.Errorf("No whitespace allowed")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	if !regexp.MustCompile(emailRegex).MatchString(info.Email) {
		return fmt.Errorf("Email format is not valid")
	}

	if len(strings.TrimSpace(info.Password)) < 6 || !regexp.MustCompile("[0-9]").MatchString(info.Password) || !regexp.MustCompile("[A-Z]").MatchString(info.Password) {
		return fmt.Errorf("Password must be 6 or more characters, include a capital letter, and include a number")
	}

	return nil
}

// Get Functions
func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Invalid id, given %s", idStr)
	}

	return id, nil
}
