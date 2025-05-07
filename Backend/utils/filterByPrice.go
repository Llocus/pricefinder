package utils

import (
	"fmt"
	"price_server/types"
	"strconv"
)

func PriceHigher(e []types.Product, price string) []types.Product {
	var result []types.Product
	FPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
	}
	for _, product := range e {
		v, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			fmt.Println(err)
		}
		if v >= FPrice {
			result = append(result, product)
		}
	}
	return result
}

func PriceLower(e []types.Product, price string) []types.Product {
	var result []types.Product
	FPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
	}
	for _, product := range e {
		v, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			fmt.Println(err)
		}
		if v <= FPrice {
			result = append(result, product)
		}
	}
	return result
}
