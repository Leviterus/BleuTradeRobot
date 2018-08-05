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

type BleuAPI struct {
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

	// Build the request
	// req, err := http.NewRequest("GET", urlS, nil)
	// if err != nil {
	// 	log.Fatal("NewRequest: ", err)
	// 	return
	// }

	// // For control over HTTP client headers,
	// // redirect policy, and other settings,
	// // create a Client
	// // A Client is an HTTP client
	// client := &http.Client{}

	// // Send the request via a client
	// // Do sends an HTTP request and
	// // returns an HTTP response
	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Fatal("Do: ", err)
	// 	return
	// }

	// // Callers should close resp.Body
	// // when done reading from it
	// // Defer the closing of the body
	// defer resp.Body.Close()

	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }

	// var record BleuAPI
	// // // fmt.Println(record)

	// // Use json.Decode for reading streams of JSON data
	// // if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
	// // 	log.Println(err)
	// // }

	// fmt.Println(record)

	// pingJSON := make(map[string][]BleuAPI)
	// erro := json.Unmarshal([]byte(resp.Body), &pingJSON)

	// if err != nil {
	// 	panic(err)
	// }

	// url := urlS

	spaceClient := http.Client{
		Timeout: time.Second * 4, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, "https://bleutrade.com/api/v2/account/getbalances?apikey=700e123416e568aeec72b5e9313d8b00&apisign=3cbd5b8a556003d210bcac4942c3e9a7f8d8cb4053950c2149c8f6d135a8e21c670a310ad696140af74ac3929bfbe059ff36c26cfc6830256aed0d63873c74a2", nil)
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

	people1 := BleuAPI{}
	fmt.Println(people1)

	jsonErr := json.Unmarshal(body, &people1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(people1)

	// m.Vars contains all variables, parsed during routing
	for _, v := range people1.Result {
		fmt.Println(v.Balance)
		m.Reply(v.Balance)
	}

}

func withdrawHandler(m *tbot.Message) {

	safeKey := url.QueryEscape(apikey)

	urlU := fmt.Sprintf("https://bleutrade.com/api/v2/account/getbalances?apikey='%s'", safeKey)

	signkey := ComputeHmac512(urlU, apisecret)
	safeSign := url.QueryEscape(signkey)

	urlS := fmt.Sprintf("https://bleutrade.com/api/v2/account/getbalances?apikey='%s'&apisign='%s'", safeKey, safeSign)

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
	var record BleuAPI

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	// m.Vars contains all variables, parsed during routing
	// 	for _, v := range record.Result {
	// 		m.Reply(FloatToString(v.Balance))
	// 	}
}

func tradeHandler(m *tbot.Message) {
	// m.Vars contains all variables, parsed during routing
	apikey = m.Vars["key"]

	fmt.Println(apikey)
	// Convert string variable to integer seconds value
	seconds, err := strconv.Atoi(apikey)
	if err != nil {
		m.Reply("Invalid API Key")
		return
	}
	m.Replyf("Timer for %d seconds started", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	m.Reply("Time out!")
}
