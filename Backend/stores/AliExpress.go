package Stores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"price_server/types"
)

func getAliExpressRequest(url string, payloadBytes []byte) string {
	data := payloadBytes
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return ""
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		return string(body)
	}
	return ""
}

type Data struct {
	Result Result `json:"result"`
}

type Result struct {
	PageInfo PageInfo `json:"pageInfo"`
	Mods     Mods     `json:"mods"`
}

type PageInfo struct {
	PageSize int `json:"pageSize"`
}

type Mods struct {
	ItemList ItemList `json:"itemList"`
}

type ItemList struct {
	Content []Content `json:"content"`
}

type Content struct {
	Image      Image      `json:"image"`
	Title      Title      `json:"title"`
	Prices     Prices     `json:"prices"`
	ProductId  string     `json:"productId"`
	Evaluation Evaluation `json:"evaluation"`
}

type Image struct {
	ImgUrl string `json:"imgUrl"`
}

type Title struct {
	DisplayTitle string `json:"displayTitle"`
}

type Prices struct {
	SalePrice     SalePrice     `json:"salePrice"`
	OriginalPrice OriginalPrice `json:"originalPrice"`
}

type SalePrice struct {
	FormattedPrice   string  `json:"formattedPrice"`
	MinPriceDiscount float64 `json:"minPriceDiscount"`
	Cent             float64 `json:"cent"`
	MinPrice         float64 `json:"minPrice"`
	Discount         float64 `json:"discount"`
}

type OriginalPrice struct {
	FormattedPrice string  `json:"formattedPrice"`
	Cent           float64 `json:"cent"`
	MinPrice       float64 `json:"minPrice"`
}

type Evaluation struct {
	StarRating float64 `json:"starRating"`
}

type AliExpressSearchResponse struct {
	Data Data `json:"data"`
}

func AliExpressSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	queryString := "https://pt.aliexpress.com/fn/search-pc/index"

	pr := ""
	if lowPrice != "" && lowPrice != "0" && highPrice != "" && highPrice != "0" {
		pr = fmt.Sprintf(`%s-%s`, lowPrice, highPrice)
	} else {
		if lowPrice != "" && lowPrice != "0" {
			pr = lowPrice
		}
		if highPrice != "" && highPrice != "0" {
			pr = fmt.Sprintf(`-%s}}`, highPrice)
		}
	}

	_ = pr

	data := map[string]interface{}{
		"pageVersion": "ff8ad60b0a0d1fbfc9e484ea303a7f44",
		"target":      "root",
		"data": map[string]interface{}{
			"page":             page,
			"g":                "y",
			"SearchText":       query,
			"selectedSwitches": "pop:true",
			"sortType":         "price_asc",
			"pr":               pr,
			"shpf_co":          "BR",
			"origin":           "y",
		},
		"eventName":  "onChange",
		"dependency": []interface{}{},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("AliExpress [0] Error: ", err)
		return productList
	}

	html := getAliExpressRequest(queryString, payloadBytes)

	var searchResponse AliExpressSearchResponse
	err = json.Unmarshal([]byte(html), &searchResponse)
	if err != nil {
		fmt.Println("AliExpress [1] Error: ", err)
		return productList
	}

	for _, result := range searchResponse.Data.Result.Mods.ItemList.Content {
		finalPrice := ""
		if result.Prices.OriginalPrice.MinPrice > 0 {
			finalPrice = fmt.Sprintf("%f", result.Prices.OriginalPrice.MinPrice)
		} else {
			finalPrice = fmt.Sprintf("%f", result.Prices.SalePrice.MinPrice)
		}
		productList = append(productList, types.Product{Store: "AliExpress", TotalPages: fmt.Sprint(searchResponse.Data.Result.PageInfo.PageSize), Image: "https:" + result.Image.ImgUrl + "_.webp", Name: result.Title.DisplayTitle, Link: fmt.Sprintf("https://pt.aliexpress.com/item/%s.html", result.ProductId), Stars: fmt.Sprintf("%f", result.Evaluation.StarRating), StarsQty: "0", Price: finalPrice})
	}
	return productList
}
