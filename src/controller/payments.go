package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Users represents a list of users
type Users []User

// User represents a user of the payment system
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Accounts  []int  `json:"accounts"`
	Transfers []int  `json:"transfers"`
	Likes     []int  `json:"likes"`
}

// Accounts -- list of account
type Accounts []Account

// Account -- represents user account type
type Account struct {
	ID            int    `json:"id"`
	User          int    `json:"user"`
	AccountNumber string `json:"accountNumber"`
	Balance       int    `json:"balance"`
}

// Transfers -- list of transfer type
type Transfers []Transfer

// Transfer -- type representing transfer of funds between origin and target users
type Transfer struct {
	ID            int        `json:"id"`
	Status        string     `json:"status"`
	OriginAccount int        `json:"originAccount"`
	TargetAccount int        `json:"targetAccount"`
	Amount        int        `json:"amount"`
	Description   string     `json:"description"`
	InitiatedAt   *time.Time `json:"initiatedAt"`
	CompletedAt   *time.Time `json:"completedAt"`
	FailedAt      *time.Time `json:"failedAt"`
}

// Likes -- list of type like
type Likes []Like

// Like -- like associated with a transfer type
type Like struct {
	ID       int `json:"id"`
	User     int `json:"user"`
	Transfer int `json:"transfer"`
}

// TransfersFeed -- list of payment transfer type
type TransfersFeed []PaymentTransfer

// PaymentTransfer -- type representing payment transfers by users
type PaymentTransfer struct {
	OriginUserName   string  `json:"originUsername"`
	TargetedUserName string  `json:"targetedUserName"`
	Amount           float64 `json:"amount"`
	Descriptions     string  `json:"description"`
	LikesCount       int     `json:"likesCount"`
}

// TransferResource -- type representing new request transfers and the response to created transfers
type TransferResource struct {
	OriginAccount int       `json:"originAccount"`
	TargetAccount int       `json:"targetAccount"`
	Amount        int       `json:"amount"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	InitiatedAt   time.Time `json:"initiatedAt"`
}

func createTransfer(payload *TransferResource) TransferResource {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://ellevest-cameron-marshall-3.glitch.me/transfers/"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	var transferResponse TransferResource
	err = json.NewDecoder(response.Body).Decode(&transferResponse)
	if err != nil {
		log.Fatal(err)
	}

	return transferResponse
}

func GetPaymentsAPIData(usersResponse *Users, transfersResponse *Transfers, likesResponse *Likes, accountsResponse *Accounts, wgcall *sync.WaitGroup) {
	go func() {
		usersHTTPResponse, err := http.Get("https://ellevest-cameron-marshall-3.glitch.me/users")
		if err != nil {
			panic(err)
		}
		defer usersHTTPResponse.Body.Close()
		err = json.NewDecoder(usersHTTPResponse.Body).Decode(&usersResponse)
		if err != nil {
			fmt.Println("error decoding json", err.Error())
		}
		wgcall.Done()
	}()
	go func() {
		accountsHTTPResponse, err := http.Get("https://ellevest-cameron-marshall-3.glitch.me/accounts")
		if err != nil {
			panic(err)
		}
		defer accountsHTTPResponse.Body.Close()
		err = json.NewDecoder(accountsHTTPResponse.Body).Decode(&accountsResponse)
		if err != nil {
			fmt.Println("error decoding json", err.Error())
		}
		wgcall.Done()
	}()

	go func() {
		transfersHTTPResponse, err := http.Get("https://ellevest-cameron-marshall-3.glitch.me/transfers")
		if err != nil {
			panic(err)
		}
		defer transfersHTTPResponse.Body.Close()
		err = json.NewDecoder(transfersHTTPResponse.Body).Decode(&transfersResponse)
		if err != nil {
			fmt.Println("error decoding json", err.Error())
		}
		wgcall.Done()
	}()
	go func() {
		likesHTTPResponse, err := http.Get("https://ellevest-cameron-marshall-3.glitch.me/likes")
		if err != nil {
			panic(err)
		}
		defer likesHTTPResponse.Body.Close()
		err = json.NewDecoder(likesHTTPResponse.Body).Decode(&likesResponse)
		if err != nil {
			println("error decoding json", err.Error())
		}
		wgcall.Done()
	}()
}

// GetPaymentsController -- handles GET requests for viewing all user transfer resources
func GetPaymentsController(w http.ResponseWriter, r *http.Request) {

	var usersResponse Users
	var accountsResponse Accounts
	var transfersResponse Transfers
	var likesResponse Likes

	wgcall := &sync.WaitGroup{}
	wgcall.Add(4)

	GetPaymentsAPIData(&usersResponse, &transfersResponse, &likesResponse, &accountsResponse, wgcall)
	wgcall.Wait()

	users := make(map[int]User)
	transfers := make(map[int]Transfer)
	likes := make(map[int]Like)
	accounts := make(map[int]Account)

	for _, user := range usersResponse {
		if _, found := users[user.ID]; !found {
			users[user.ID] = user
		}
	}
	for _, account := range accountsResponse {
		if _, found := accounts[account.ID]; !found {
			accounts[account.ID] = account
		}
	}
	for _, like := range likesResponse {
		if _, found := likes[like.ID]; !found {
			likes[like.ID] = like
		}
	}
	for _, transfer := range transfersResponse {
		if _, found := transfers[transfer.ID]; !found {
			transfers[transfer.ID] = transfer
		}
	}

	var transferFeed TransfersFeed

	for _, transfer := range transfersResponse {
		if len(transfer.Status) > 0 && transfer.Status != "failed" {
			originUserName := strings.Join([]string{users[accounts[transfer.OriginAccount].User].FirstName, users[accounts[transfer.OriginAccount].User].LastName}, " ")
			targetedUserName := strings.Join([]string{users[accounts[transfer.TargetAccount].User].FirstName, users[accounts[transfer.TargetAccount].User].LastName}, " ")
			amount := float64(transfer.Amount)
			description := transfer.Description
			likesCount := 0
			for _, like := range likes {
				if like.Transfer == transfer.ID {
					likesCount++
				}
			}
			paymentTransfer := PaymentTransfer{
				OriginUserName:   originUserName,
				TargetedUserName: targetedUserName,
				Amount:           amount,
				Descriptions:     description,
				LikesCount:       likesCount,
			}

			transferFeed = append(transferFeed, paymentTransfer)
		}
	}

	response, err := json.Marshal(transferFeed)
	if err != nil {
		fmt.Println("error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// CreateTransferController -- handles POST requests for creating new transfers
func CreateTransferController(w http.ResponseWriter, r *http.Request) {
	var payload TransferResource
	json.NewDecoder(r.Body).Decode(&payload)

	payload.InitiatedAt = time.Now()
	url := fmt.Sprintf("https://ellevest-cameron-marshall-3.glitch.me/accounts/%d", payload.OriginAccount)
	originAccountResponse, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer originAccountResponse.Body.Close()

	var account Account
	json.NewDecoder(originAccountResponse.Body).Decode(&account)

	if account.Balance >= payload.Amount {
		transferResponse := createTransfer(&payload)
		response, err := json.Marshal(transferResponse)
		if err != nil {
			fmt.Println("error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusFailedDependency)
			w.Write([]byte("Failed to create transfer"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not enough funds for transfer"))
	}
}
