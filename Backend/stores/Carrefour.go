package Stores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"price_server/types"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getCarrefourRequest(url string) string {
	data := []byte(``)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Carrefour [0] Error: ", err)
		return ""
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Carrefour [1] Error: ", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)

	}
	return ""
}

type ItemListElement struct {
	Item Item `json:"item"`
}

type Item struct {
	Id    string `json:"@id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Offer Offer  `json:"offers"`
}

type Offer struct {
	LowPrice  float64  `json:"lowPrice"`
	HighPrice float64  `json:"highPrice"`
	Offers    []Offers `json:"offers"`
}

type Offers struct {
	Price float64 `json:"price"`
}

type CarrefourSearchResponse struct {
	ItemListElement []ItemListElement `json:"itemListElement"`
}

func CarrefourSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	itensPerPage := 60
	queryString := fmt.Sprintf("https://www.carrefour.com.br/busca/%s?count=%d&sort=price_asc", query, itensPerPage)

	if page > 0 {
		queryString += fmt.Sprintf("&page=%d", page)
	}
	if lowPrice != "" && lowPrice != "0" && highPrice != "" && highPrice != "0" {
		queryString += fmt.Sprintf(`&c_price=%s:%s`, lowPrice, highPrice)
	} else {
		if lowPrice != "" && lowPrice != "0" {
			queryString += fmt.Sprintf(`&c_price=%s:1000000`, lowPrice)
		}
		if highPrice != "" && highPrice != "0" {
			queryString += fmt.Sprintf(`&c_price=0:%s`, highPrice)
		}
	}

	html := getCarrefourRequest(queryString)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Carrefour [2] Error: ", err)
	}
	htmljson := string(doc.Find(`script[type="application/ld+json"]`).First().Text())

	var searchResponse CarrefourSearchResponse
	err = json.Unmarshal([]byte(htmljson), &searchResponse)
	if err != nil {
		fmt.Println("Carrefour [3] Error: ", err)
		return productList
	}

	totalPages := 1
	totalProducts := 0
	stringProducts := doc.Find(`div[class*="totalProducts"] > span`).First().Text()
	reStringProducts := regexp.MustCompile("[0-9]+")
	matcheStringProduct := reStringProducts.FindAllString(stringProducts, -1)
	if len(matcheStringProduct) > 0 {
		totalProducts, _ = strconv.Atoi(matcheStringProduct[0])
	}
	if (itensPerPage * page) > totalProducts {
		totalPages = page + 1
	}
	for _, result := range searchResponse.ItemListElement {
		if result.Item.Offer.LowPrice > 0 {
			productList = append(productList, types.Product{Store: "Carrefour", TotalPages: fmt.Sprint(totalPages), Image: result.Item.Image, Name: result.Item.Name, Link: result.Item.Id, Stars: "-1", StarsQty: "0", Price: fmt.Sprintf("%f", result.Item.Offer.LowPrice)})
		}
	}
	return productList
}
