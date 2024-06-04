package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Smtp struct {
		IsDeliverable bool `json:"is_deliverable"`
	} `json:"smtp"`
}

func main() {
	fileName := flag.String("file", "", "input file")
	flag.Parse() // Parse the command-line arguments

	if *fileName == "" {
		log.Fatal("Please specify a csv file using -file flag")
	}
	file, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		email := record[0]
		checkEmail(email)
	}
}

func checkEmail(email string) {
	url := "https://reacher.fatlab.io/v0/check_email"
	var jsonData = []byte(fmt.Sprintf(`{
		"to_email": "%s"
	}`, email))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-reacher-secret", "aaa")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response Response
	json.Unmarshal(body, &response)

	if response.Smtp.IsDeliverable {
		fmt.Println(email)
	}
}
