package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Post struct {
	ID      string  `json:"id"`
	Receipt Receipt `json:"receipt"`
	Point   int64   `json:"point"`
}

type postResponse struct {
	ID string `json:"id"`
}

type getResponse struct {
	Points int64 `json:"points"`
}

var (
	posts   = make(map[string]Post)
	postsMu sync.Mutex
)

func main() {
	http.HandleFunc("/receipts/", postsHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.URL.Path[len("/receipts/") : len(r.URL.Path)-len("/points")]
		handleGetPost(w, r, id)
	case "POST":
		handlePostPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostPost(w http.ResponseWriter, r *http.Request) {
	var p Post
	var reqData Receipt
	var res postResponse
	var validString = regexp.MustCompile("^[\\w\\s\\-]+$")
	var validPrice = regexp.MustCompile("^\\d+\\.\\d{2}$")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &reqData); err != nil || reqData.Items == nil || len(reqData.Items) == 0 || len(reqData.Retailer) == 0 || len(reqData.PurchaseDate) == 0 || len(reqData.PurchaseTime) == 0 || len(reqData.Total) == 0 {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}
	date, err := time.Parse(time.DateOnly, reqData.PurchaseDate)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	time, err := time.Parse("15:04", reqData.PurchaseTime)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	if isMatch, err := regexp.MatchString("^[\\w\\s\\-&]+$", reqData.Retailer); err != nil || !isMatch {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	if !validPrice.MatchString(reqData.Total) {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	for _, e := range reqData.Items {
		if len(e.Price) == 0 || len(e.ShortDescription) == 0 {
			http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
			return
		} else if !validString.MatchString(e.ShortDescription) {
			http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
			return
		} else if !validPrice.MatchString(e.Price) {
			http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
			return
		}
	}
	postsMu.Lock()
	defer postsMu.Unlock()

	p.ID = uuid.New().String()
	p.Point = int64(calculatePoints(reqData.Retailer, reqData.Items, date, time, reqData.Total))
	p.Receipt = reqData

	posts[p.ID] = p
	res.ID = p.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func handleGetPost(w http.ResponseWriter, r *http.Request, id string) {
	postsMu.Lock()
	defer postsMu.Unlock()

	p, ok := posts[id]
	if !ok {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}
	var res getResponse
	res.Points = p.Point
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func calculatePoints(retailer string, items []Item, date time.Time, hour time.Time, total string) int64 {
	var point int64 = 0

	//One point for every alphanumeric character in the retailer name.
	for _, char := range retailer {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
			point++
		}
	}

	//50 points if the total is a round dollar amount with no cents.
	cents, _ := strconv.Atoi(total[len(total)-2:])
	if cents == 0 {
		point += 50
	}

	//25 points if the total is a multiple of 0.25
	if cents%25 == 0 {
		point += 25
	}

	//5 points for every two items on the receipt.
	point += int64(len(items)) / 2 * 5

	//If the trimmed length of the item description is a multiple of 3,
	//multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(item.Price, 64)
			point += int64(math.Ceil(itemPrice * 0.2))
		}
	}

	//6 points if the day in the purchase date is odd.
	if date.Day()%2 != 0 {
		point += 6
	}

	//10 points if the time of purchase is after 2:00pm and before 4:00pm.
	afterTwo, _ := time.Parse("15:04", "14:00")
	beforeFour, _ := time.Parse("15:04", "16:00")

	if hour.After(afterTwo) && hour.Before(beforeFour) {
		point += 10
	}
	return point
}
