package Stores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"price_server/types"
	"price_server/utils"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getMLRequest(url string) string {
	data := []byte(``)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("ML [0] Error: ", err)
		return ""
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("ML [1] Error: ", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			fmt.Println("ML [2] Error: ", err)
		}
		docRes := doc.Find("#__PRELOADED_STATE__").Text()

		return docRes
	}
	return ""
}

type MLPageState struct {
	InitialState MLInitialState `json:"initialState"`
}

type MLPagination struct {
	PageCount    int `json:"page_count"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	SelectedPage int `json:"selected_page"`
}

type MLAggregateRating struct {
	RatingCount float64 `json:"rating_count"`
	RatingValue float64 `json:"rating_value"`
}

type MLItemOffered struct {
	Price float64 `json:"price"`
	Url   string  `json:"url"`
}

type MLProduct struct {
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	ItemOffered     MLItemOffered     `json:"item_offered"`
	AggregateRating MLAggregateRating `json:"aggregate_rating"`
}

type MLSchema struct {
	Products []MLProduct `json:"product_list"`
}

type MLSeo struct {
	Schema MLSchema `json:"schema"`
}

type MLInitialState struct {
	Query      string       `json:"query"`
	Pagination MLPagination `json:"pagination"`
	Results    []MLResult   `json:"results"`
	Seo        MLSeo        `json:"seo"`
}

type MLPicturesStack struct {
	Retina string `json:"retina"`
}

type MLPictures struct {
	Stack MLPicturesStack `json:"stack"`
}

type MLResult struct {
	ID        string      `json:"id"`
	Title     string      `json:"title,omitempty"`
	Price     MLPriceInfo `json:"price"`
	Pictures  MLPictures  `json:"pictures"`
	Permalink string      `json:"permalink"`
	Reviews   MLReviews   `json:"reviews"`
	IsAd      bool        `json:"is_ad"`
	AdLabel   string      `json:"ad_label"`
}

type MLReviews struct {
	RatingAverage float64 `json:"rating_average"`
	Total         int     `json:"total"`
}

type MLPriceInfo struct {
	Amount     float64 `json:"amount"`
	CurrencyID string  `json:"currency_id"`
}

type MLSearchResponse struct {
	PageState MLPageState `json:"pageState"`
}

func MLSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	queryString := fmt.Sprintf("https://lista.mercadolivre.com.br/%s", strings.Replace(query, " ", "-", -1))
	queryString += "_OrderId_PRICE"
	if page > 1 {
		queryString += fmt.Sprintf("_Desde_%d", ((page-1)*50)+1)
	}
	if lowPrice != "" && lowPrice != "0" && highPrice != "" && highPrice != "0" {
		queryString += fmt.Sprintf("_PriceRange_%s-%s", lowPrice, highPrice)
	} else {
		if lowPrice != "" && lowPrice != "0" {
			queryString += fmt.Sprintf("_PriceRange_%s-0", lowPrice)
		}
		if highPrice != "" && highPrice != "0" {
			queryString += fmt.Sprintf("_PriceRange_0-%s", highPrice)
		}
	}
	queryString += "_NoIndex_True"

	html := getMLRequest(queryString)

	var searchResponse MLSearchResponse
	err := json.Unmarshal([]byte(html), &searchResponse)
	if err != nil {
		fmt.Println("ML [3] Error: ", err)
	}

	for _, result := range searchResponse.PageState.InitialState.Seo.Schema.Products {
		if result.ItemOffered.Price > 0 && result.Name != "" /*&& result.Reviews.RatingAverage >= 3*/ {
			parsedURL, err := url.Parse(result.ItemOffered.Url)
			if err != nil {
				fmt.Println("ML [4] Error: ", err)
			}
			params := parsedURL.Query()
			delete(params, "position")
			delete(params, "tracking_id")
			parsedURL.RawQuery = params.Encode()
			productList = append(productList, types.Product{Store: "MercadoLivre", TotalPages: fmt.Sprintf("%d", searchResponse.PageState.InitialState.Pagination.LastPage), Image: result.Image, Name: result.Name, Link: parsedURL.String(), Stars: fmt.Sprintf("%f", result.AggregateRating.RatingValue), StarsQty: fmt.Sprintf("%f", result.AggregateRating.RatingCount), Price: strings.TrimSpace(fmt.Sprintf("%f", result.ItemOffered.Price))})
		}
	}
	sort.Slice(productList, utils.OrderByPrice(productList))
	return productList
}
