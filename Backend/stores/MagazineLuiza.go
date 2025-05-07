package Stores

import (
	"bytes"
	"crypto/tls"
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

func getMagazineLuizaRequest(url string) string {
	data := []byte(``)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("MagazineLuiza [0] Error: ", err)
		return ""
	}
	request.Header.Set("Host", "www.magazineluiza.com.br")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("MagazineLuiza [1] Error: ", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)

	}
	return ""
}

type Graph struct {
	Sku             string             `json:"sku"`
	Name            string             `json:"name"`
	Image           string             `json:"image"`
	AggregateRating AggregateRating    `json:"aggregateRating"`
	Offer           MagazineLuizaOffer `json:"offers"`
}

type AggregateRating struct {
	RatingValue string `json:"ratingValue"`
	ReviewCount string `json:"reviewCount"`
}

type MagazineLuizaOffer struct {
	Price string `json:"price"`
	Url   string `json:"url"`
}

type MagazineLuizaOffers struct {
	Price float64 `json:"price"`
}

type MagazineLuizaSearchResponse struct {
	Graph []Graph `json:"@graph"`
}

func MagazineLuizaSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	itensPerPage := 60
	queryString := strings.ReplaceAll(fmt.Sprintf("https://www.magazineluiza.com.br/busca/%s/?sortOrientation=asc&sortType=price", strings.ToLower(query)), " ", "+")

	if page > 0 {
		queryString += fmt.Sprintf("&page=%d", page)
	}
	if lowPrice != "" && lowPrice != "0" && highPrice != "" && highPrice != "0" {
		queryString += fmt.Sprintf(`filters=price---%s:%s`, lowPrice, highPrice)
	} else {
		if lowPrice != "" && lowPrice != "0" {
			queryString += fmt.Sprintf(`filters=price---%s:1000000`, lowPrice)
		}
		if highPrice != "" && highPrice != "0" {
			queryString += fmt.Sprintf(`filters=price---0:%s`, highPrice)
		}
	}

	html := getMagazineLuizaRequest(queryString)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("MagazineLuiza [2] Error: ", err)
	}

	htmljson := ""
	doc.Find(`script[type="application/ld+json"]`).Each(func(j int, s2 *goquery.Selection) {
		if strings.Contains(s2.Text(), "offers") {
			htmljson = string(s2.Text())
		}
	})

	var searchResponse MagazineLuizaSearchResponse
	err = json.Unmarshal([]byte(htmljson), &searchResponse)
	if err != nil {
		fmt.Println("MagazineLuiza [3] Error: ", err)
		return productList
	}

	totalPages := 1
	totalProducts := 0
	stringProducts := doc.Find(`div[data-testid="mod-q"] p`).First().Text()
	reStringProducts := regexp.MustCompile("[0-9]+")
	matcheStringProduct := reStringProducts.FindAllString(stringProducts, -1)
	if len(matcheStringProduct) > 0 {
		totalProducts, _ = strconv.Atoi(matcheStringProduct[0])
	}
	if (itensPerPage * page) > totalProducts {
		totalPages = page + 1
	}
	for _, result := range searchResponse.Graph {
		price, _ := strconv.ParseFloat(result.Offer.Price, 64)
		if price > 0 {
			productList = append(productList, types.Product{Store: "MagazineLuiza", TotalPages: fmt.Sprint(totalPages), Image: result.Image, Name: result.Name, Link: result.Offer.Url, Stars: result.AggregateRating.RatingValue, StarsQty: result.AggregateRating.ReviewCount, Price: result.Offer.Price})
		}
	}
	return productList
}
