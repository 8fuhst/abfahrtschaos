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

type departure_list_body struct {
	Station string `json:"station"`
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getSignature(request_body []byte) string {
	pass := os.Getenv("API_PASS")
	pass_bytes := []byte(pass) //key

	//bytes_to_hash := append(pass_bytes, request_body...)
	//hashed := hmac.New(sha1.New, bytes_to_hash)
	hashed := hmac.New(sha1.New, pass_bytes)
	_, err := hashed.Write(request_body)
	if err != nil {
		log.Fatalf("Error during signature generation")
	}
	signature := base64.StdEncoding.EncodeToString(hashed.Sum(nil))
	return signature
}

func main() {
	username := getEnvVariable("API_USERNAME")
	fmt.Print("Starting HTTP Client...")
	client := &http.Client{}

	d := departure_list_body{"Altona"}
	body, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error while creating JSON")
	}
	body_buffer := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", "http://gti.geofox.de/gti/public/departureList", body_buffer)

	if err != nil {
		log.Fatal(err)
	}

	signature := getSignature(body)
	fmt.Print(signature)

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

	resp_body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(resp_body)
	fmt.Print(bodyString)
}
