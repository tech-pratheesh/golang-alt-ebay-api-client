package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
	"github.com/mitchellh/mapstructure"
)

var EBAY_SOLD_PRICE_QUERY = `
  query EbayItemPrice($input: EbayItemPriceInput!) {
    ebayItemPrice(input: $input) {
      ebayItemPrice {
        itemId
        price
      }
    }
  }`
var ALT_GRAPHQL_ENDPOINT = "https://alt-platform-server.production.internal.onlyalt.com/graphql/"

type PriceQueryResponseModel struct {
	EbayItemPrice map[string]EbayItemPrice `json:"ebayItemPrice"`
}

type EbayItemPrice struct {
	ItemID string `json:"itemId"`
	Price  string `json:"price"`
}

func main() {
	// define flag
	flagItemId := flag.String("item-id", "", "")
	flagListingType := flag.String("listing-type", "BEST_OFFER_ACCEPTED", "")
	flag.Parse()

	if *flagItemId == "" {
		log.Fatalf("you should provide an item id")
		return
	}
	// initiate client
	graphqlClient := graphql.NewClient(ALT_GRAPHQL_ENDPOINT)

	// make a request
	graphqlRequest := graphql.NewRequest(EBAY_SOLD_PRICE_QUERY)

	// set variable
	graphqlRequest.Var("input", map[string]string{
		"itemId":      *flagItemId,
		"listingType": *flagListingType,
	})

	// set content type
	graphqlRequest.Header.Set("Content-Type", "application/json")

	var graphqlResponse interface{}
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	responseData := graphqlResponse.(map[string]interface{})

	data := PriceQueryResponseModel{}
	mapstructure.Decode(responseData, &data)

	fmt.Println("ItemID : ", data.EbayItemPrice["ebayItemPrice"].ItemID)
	fmt.Println("Price  : ", data.EbayItemPrice["ebayItemPrice"].Price)
}
