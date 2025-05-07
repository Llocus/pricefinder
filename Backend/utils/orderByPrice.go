package utils

import (
	"fmt"
	"price_server/types"
	"strconv"
)

func OrderByPrice(e []types.Product) func(i, j int) bool {
	return func(i, j int) bool {
		v1, err := strconv.ParseFloat(e[i].Price, 64)
		if err != nil {
			fmt.Println(err)
		}
		v2, err := strconv.ParseFloat(e[j].Price, 64)
		if err != nil {
			fmt.Println(err)
		}
		return v1 < v2
	}
}
