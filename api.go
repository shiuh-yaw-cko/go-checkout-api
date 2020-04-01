package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var sKey = "sk_test_79ab19cc-5c16-4b81-8110-31666040bb6a"
var pKey = "pk_test_291004cb-8a16-44c3-8c64-37ec45b47cd4"
var host = "https://api.sandbox.checkout.com/"
var contentType = "application/json"

// PaymentInfo ...
type PaymentInfo struct {
	ID              string   `json:"id"`
	ActionID        string   `json:"action_id"`
	Amount          int      `json:"amount"`
	Currency        string   `json:"currency"`
	Approved        bool     `json:"approved"`
	Status          string   `json:"status"`
	AuthCode        string   `json:"auth_code"`
	ResponseCode    string   `json:"response_code"`
	ResponseSummary string   `json:"response_summary"`
	ProcessedOn     string   `json:"processed_on"`
	Reference       string   `json:"reference"`
	Risk            Risk     `json:"risk"`
	Source          Source   `json:"source"`
	Customer        Customer `json:"customer"`
	Links           Links    `json:"_links"`
}

// Risk ...
type Risk struct {
	Flagged bool `json:"flagged"`
}

// Source ...
type Source struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	ExpiryMonth   int    `json:"expiry_month"`
	ExpiryYear    int    `json:"expiry_year"`
	Scheme        string `json:"scheme"`
	Last4         string `json:"last4"`
	Fingerprint   string `json:"fingerprint"`
	Bin           string `json:"bin"`
	CardType      string `json:"card_type"`
	CardCategory  string `json:"card_category"`
	Issuer        string `json:"issuer"`
	IssuerCountry string `json:"issuer_country"`
	ProductID     string `json:"product_id"`
	ProductType   string `json:"product_type"`
}

// Customer ...
type Customer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Links ...
type Links struct {
	Current     URL `json:"self"`
	RedirectURL URL `json:"redirect"`
}

// URL ...
type URL struct {
	URLString string `json:"href"`
}

// CreditCard User credit card information
type CreditCard struct {
	Type        string `json:"type"`
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	Name        string `json:"name"`
	CVV         string `json:"cvv"`
}

// Token ...
type Token struct {
	Type          string `json:"type"`
	Token         string `json:"token"`
	ExpiresOn     string `json:"expires_on"`
	ExpiryMonth   int    `json:"expiry_month"`
	ExpiryYear    int    `json:"expiry_year"`
	Scheme        string `json:"scheme"`
	Last4         string `json:"last4"`
	Bin           string `json:"bin"`
	CardType      string `json:"card_type"`
	CardCategory  string `json:"card_category"`
	Issuer        string `json:"issuer"`
	IssuerCountry string `json:"issuer_country"`
	ProductID     string `json:"product_id"`
	ProductType   string `json:"product_type"`
	Name          string `json:"name"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Shiuh Yaw Checkout's API")
}

func pay(w http.ResponseWriter, r *http.Request) {

	url := host + "tokens"
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		log.Printf("http.NewRequest() error: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", pKey)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}
	var token Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		fmt.Printf("Unmarshal : %v\n", err)
		return
	}
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(token)
	var creditCard CreditCard
	requestPayment(w, r, token, creditCard)
}

func requestPayment(w http.ResponseWriter, r *http.Request, token Token, card CreditCard) {

	body := map[string]interface{}{
		"source": map[string]string{
			"type":  "token",
			"token": token.Token,
		},
		"amount":    "2500",
		"currency":  "GBP",
		"reference": "Test Order",
		// "3ds": map[string]bool{
		// 	"enabled":     true,
		// 	"attempt_n3d": true,
		// },
		"customer": map[string]string{
			"name": card.Name,
		},
	}
	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}

	paymentURL := host + "payments"
	req, err := http.NewRequest(r.Method, paymentURL, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Printf("http.NewRequest() error: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", sKey)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	var response PaymentInfo
	err = json.Unmarshal(data, &response)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func getPaymentInfo(w http.ResponseWriter, r *http.Request) {

	paymentID := mux.Vars(r)["payment_id"]
	if len(paymentID) < 0 {
		fmt.Fprintf(w, "Payment ID invalid")
		return
	}
	url := host + "payments/" + paymentID
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		log.Printf("http.NewRequest() error: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", sKey)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}
	var info PaymentInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		fmt.Printf("Unmarshal : %v\n", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/pay", pay).Methods("POST")
	router.HandleFunc("/info/{payment_id}", getPaymentInfo).Methods("GET")
	fmt.Println("Listening")
	log.Fatal(http.ListenAndServe(":8081", router))
}
