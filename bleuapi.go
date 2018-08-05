package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yanzay/tbot"
)

type BalanceAPI struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Result  []struct {
		Currency      string `json:"Currency"`
		Balance       string `json:"Balance"`
		Available     string `json:"Available"`
		Pending       string `json:"Pending"`
		CryptoAddress string `json:"CryptoAddress"`
		IsActive      string `json:"IsActive"`
	} `json:"result"`
}

type TradeAPI struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Orderid string `json:"orderid"`
	} `json:"result"`
}

type WithdrawAPI struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func ComputeHmac512(message string, key string) string {

	keyB := []byte(key)

	sig := hmac.New(sha512.New, keyB)
	sig.Write([]byte(message))

	return hex.EncodeToString(sig.Sum(nil))

	// key := []byte(secret)
	// h := hmac.New(sha512.New, key)
	// h.Write([]byte(message))
	// return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func balanceHandler(m *tbot.Message) {

	safeKey := url.QueryEscape(apikey)
	fmt.Println(safeKey)

	urlU := fmt.Sprintf("https://bleutrade.com/api/v2/account/getbalances?apikey=%s", safeKey)
	fmt.Println(urlU)

	signkey := ComputeHmac512(urlU, apisecret)
	fmt.Println(signkey)

	safeSign := url.QueryEscape(signkey)
	fmt.Println(safeSign)

	urlS := fmt.Sprintf("https://bleutrade.com/api/v2/account/getbalances?apikey=%s&apisign=%s", safeKey, safeSign)
	fmt.Println(urlS)

	spaceClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, urlS, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	people1 := BalanceAPI{}
	fmt.Println(people1)

	jsonErr := json.Unmarshal(body, &people1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(people1)

	for _, v := range people1.Result {
		if v.Balance != "0.00000000" {
			value := fmt.Sprintf("%s: %s", v.Currency, v.Balance)
			m.Reply(value)
		}

	}

}

func withdrawHandler(m *tbot.Message) {

	address := m.Vars["adress"]
	quantity := m.Vars["quantity"]
	currency := m.Vars["currency"]
	safeKey := url.QueryEscape(apikey)

	urlU := fmt.Sprintf("https://bleutrade.com/api/v2/account/withdraw?apikey=%s", safeKey)

	signkey := ComputeHmac512(urlU, apisecret)
	safeSign := url.QueryEscape(signkey)

	urlS := fmt.Sprintf("https://bleutrade.com/api/v2/account/withdraw?apikey=%s&apisign=%s&address=%s&quantity=%s&currency=%s", safeKey, safeSign, address, quantity, currency)

	// Build the request
	req, err := http.NewRequest("GET", urlS, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record WithdrawAPI

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	if record.Success == "true" {
		m.Reply("Withdraw made with success!!")
	} else {
		m.Reply("Withdraw has failed.")
	}

	// m.Vars contains all variables, parsed during routing
	// 	for _, v := range record.Result {
	// 		m.Reply(FloatToString(v.Balance))
	// 	}
}

func tradeBuyHandler(m *tbot.Message) {

	market := m.Vars["market"]
	rate := m.Vars["rate"]
	quantity := m.Vars["quantity"]
	safeKey := url.QueryEscape(apikey)

	urlU := fmt.Sprintf("https://bleutrade.com/api/v2/market/buylimit?apikey=%s", safeKey)

	signkey := ComputeHmac512(urlU, apisecret)
	safeSign := url.QueryEscape(signkey)

	urlS := fmt.Sprintf("https://bleutrade.com/api/v2/market/buylimit?apikey=%s&apisign=%s&market=%s&rate=%s&quantity=%s", safeKey, safeSign, market, rate, quantity)

	// Build the request
	req, err := http.NewRequest("GET", urlS, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record TradeAPI

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	if record.Success == "true" {
		m.Reply("Buy ORDER made with sucess!!")
	} else {
		m.Reply("Buy ORDER has failed.")
	}

	// m.Vars contains all variables, parsed during routing
	// 	for _, v := range record.Result {
	// 		m.Reply(FloatToString(v.Balance))
	// 	}
}

func tradeSellHandler(m *tbot.Message) {

	market := m.Vars["market"]
	rate := m.Vars["rate"]
	quantity := m.Vars["quantity"]
	safeKey := url.QueryEscape(apikey)

	urlU := fmt.Sprintf("https://bleutrade.com/api/v2/market/selllimit?apikey=%s", safeKey)

	signkey := ComputeHmac512(urlU, apisecret)
	safeSign := url.QueryEscape(signkey)

	urlS := fmt.Sprintf("https://bleutrade.com/api/v2/market/selllimit?apikey=%s&apisign=%s&market=%s&rate=%s&quantity=%s", safeKey, safeSign, market, rate, quantity)

	// Build the request
	req, err := http.NewRequest("GET", urlS, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record TradeAPI

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	if record.Success == "true" {
		m.Reply("Sell ORDER made with sucess!!")
	} else {
		m.Reply("Sell ORDER has failed.")
	}

	// m.Vars contains all variables, parsed during routing
	// 	for _, v := range record.Result {
	// 		m.Reply(FloatToString(v.Balance))
	// 	}
}
