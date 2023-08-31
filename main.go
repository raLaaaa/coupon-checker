package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {

	s := gocron.NewScheduler(time.UTC)

	fmt.Println("Starting routine")

	_, err := s.Every(10).Seconds().Do(task)

	if err != nil {
		print(err)
	}

	_, err = os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	s.StartBlocking()
}

func task() {

	t := fmt.Sprint(time.Now().Nanosecond())

	resp, err := http.Get("https://gonzalesbremen.de/api/shop/d176b2a7-be82-11e9-80dd-00163e41c820/validate_coupon?code=" + t[:6])
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString("Checking Code: " + t[:6] + "\n"); err != nil {
		panic(err)
	}

	if strings.Contains(string(body), "404") {
		if _, err = f.WriteString("Fail\n"); err != nil {
			panic(err)
		}
	} else {
		if _, err = f.WriteString("--------------SUCCESS----------------\n"); err != nil {
			panic(err)
		}
	}

	if _, err = f.WriteString("--------------\n"); err != nil {
		panic(err)
	}
}
