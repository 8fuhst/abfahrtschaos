package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type cnRequest struct {
	TheName     sdName `json:"theName"`
	MaxList     int    `json:"maxList"`
	MaxDistance int    `json:"maxDistance"`
}

type dlRequest struct {
	Station       sdName `json:station`
	Time          string `json:time`
	MaxList       int    `json:maxList`
	MaxTimeOffset int    `json:maxTimeOffset`
}

type sdName struct {
	Name string `json:"name"`
	City string `json:"city"`
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getSignature(requestBody []byte) string {
	pass := getEnvVariable("API_PASS")
	fmt.Printf("Password: %s\n", pass)
	pass_bytes := []byte(pass) //key

	//bytes_to_hash := append(pass_bytes, request_body...)
	//hashed := hmac.New(sha1.New, bytes_to_hash)
	hashed := hmac.New(sha1.New, pass_bytes)
	_, err := hashed.Write(requestBody)
	if err != nil {
		log.Fatalf("Error during signature generation")
	}
	signature := base64.StdEncoding.EncodeToString(hashed.Sum(nil))
	return signature
}

func execRequest(client *http.Client, signature string, req *http.Request) *http.Response {
	username := getEnvVariable("API_USERNAME")

	req.Header.Add("geofox-auth-user", username)
	req.Header.Add("geofox-auth-signature", signature)
	req.Header.Add("version", "57")
	req.Header.Add("X-Platform", "web")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func requestCheckName(stationName string, cityName string, maxList int, maxDistance int, client *http.Client) *http.Response {
	url := "http://gti.geofox.de/gti/public/checkName"
	sdname := sdName{Name: stationName, City: cityName}
	body := cnRequest{TheName: sdname, MaxList: maxList, MaxDistance: maxDistance}

	body_json, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Error while creating JSON for checkName")
	}

	body_buffer := bytes.NewBuffer(body_json)
	req, err := http.NewRequest("POST", url, body_buffer)

	if err != nil {
		log.Fatal(err)
	}

	signature := getSignature(body_json)
	fmt.Printf("Signature for checkName request: %s\n", signature)

	return execRequest(client, signature, req)
}

func requestDepartureList(stationName string, cityName string, maxList int, maxTimeOffset int, client *http.Client) *http.Response {
	url := "http://gti.geofox.de/gti/public/departureList"
	sdname := sdName{Name: stationName, City: cityName}
	body := dlRequest{Station: sdname, MaxList: maxList, MaxTimeOffset: maxTimeOffset}

	body_json, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Error while creating JSON for departureList")
	}

	body_buffer := bytes.NewBuffer(body_json)
	req, err := http.NewRequest("POST", url, body_buffer)

	if err != nil {
		log.Fatal(err)
	}

	signature := getSignature(body_json)
	fmt.Printf("Signature for departureList request: %s\n", signature)

	return execRequest(client, signature, req)
}

func main() {
	fmt.Print("Starting HTTP Client...\n")
	client := &http.Client{}

	resp := requestCheckName("Altona", "Hamburg", 1, 10, client)

	resp_body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(resp_body)
	fmt.Print(bodyString)
}
