package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Ticker struct {
	Key TickerData `json:"ticker"`
}

type TickerData struct {
	High string `json:"high"`
	Low  string `json:"low"`
}

func main() {
	err := godotenv.Load(".env")
	check_err(err)

	fetch_info()
	send_info()
}

func fetch_info() {
	url := os.Getenv("INDODAX_URL")
	last_buy_price, err := strconv.ParseFloat(os.Getenv("LAST_BUY_PRICE"), 64)
	check_err(err)
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	check_err(err)
	res, getErr := httpClient.Do(req)
	check_err(getErr)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	check_err(readErr)

	ticker1 := Ticker{}
	jsonErr := json.Unmarshal(body, &ticker1)
	check_err(jsonErr)
	high, err := strconv.ParseFloat(ticker1.Key.High, 64)
	check_err(err)
	low, err := strconv.ParseFloat(ticker1.Key.Low, 64)
	check_err(err)
	avg := float64((high + low) / 2)
	profit := (avg - last_buy_price) / last_buy_price * 100

	result := fmt.Sprintf("profit: %.2f\nbtc_idr: %.0f\nlast_buy_price: %.0f\ntime_stamp: %v", profit, avg, last_buy_price, time.Now().Format("2006-01-02 3:4:5 PM"))
	write_file(result)
}

func send_info() {
	info_byte, err := ioutil.ReadFile("savefile")
	check_err(err)

	flock_token := os.Getenv("FLOCK_TOKEN")

	params := url.Values{}
	params.Add("text", string(info_byte))

	resp, err := http.PostForm("https://api.flock.com/hooks/sendMessage/"+flock_token, params)
	check_err(err)
	fmt.Print(resp)
	defer resp.Body.Close()
}

func check_err(e error) {
	if e != nil {
		panic(e)
	}
}

func write_file(data string) {
	d1 := []byte(data)
	err := ioutil.WriteFile("savefile", d1, 0644)
	check_err(err)
}
