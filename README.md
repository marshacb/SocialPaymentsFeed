#Social Payments Server
This project allows users to view all non failed payment transfers as well as create new transfers from from 1 user to another.

It supports the following request(s):

```GET /v1/payments/```
```POST /v1/transfers/```

## Prerequisites

Be sure to have Go installed locally.

## Installation
After saving root folder locally cd into root and install dependencies:

```
cd SocialPaymentsFeed 

#Install the dependencies with go get
go get -d ./...
```

## Running the server

From the root project directory run

```go run main.go```

## Example

```http://localhost:8080/v1/payments/```
