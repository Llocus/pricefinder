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

func getExtraRequest(url string) string {
	data := []byte(``)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Extra [0] Error: ", err)
		return ""
	}
	request.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,/;q=0.8")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("Content-Type", "text/html; charset=utf-8")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Mobile Safari/537.36")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Extra [1] Error: ", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)

	}
	return ""
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	InitialState InitialState `json:"initialState"`
}

type InitialState struct {
	Search Search `json:"search"`
}

type Search struct {
	Results Results `json:"results"`
}

type Results struct {
	Size     int        `json:"size"`
	Products []Products `json:"products"`
}

type Products struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Image       string  `json:"image"`
	Rating      float64 `json:"rating"`
	RatingCount float64 `json:"ratingCount"`
	Href        string  `json:"href"`
	CId         string  `json:"cId"`
	IdSku       int     `json:"idSku"`
}

type RuntimeConfig struct {
	ResultsPerPage string `json:"RESULTS_PER_PAGE"`
	PriceApiKey    string `json:"PRICE_API_KEY"`
}

type SearchResponse struct {
	Props         Props         `json:"props"`
	RuntimeConfig RuntimeConfig `json:"runtimeConfig"`
}

/*	*/

type PrecoProdutos struct {
	PrecoVenda PrecoVenda `json:"PrecoVenda"`
}

type PrecoVenda struct {
	IdSku            int     `json:"IdSku"`
	IdProduto        int     `json:"IdProduto"`
	PrecoDe          float64 `json:"PrecoDe"`
	Preco            float64 `json:"Preco"`
	PrecoSemDesconto float64 `json:"PrecoSemDesconto"`
	NumeroParcelas   int     `json:"NumeroParcelas"`
}

type SearchResponsePrices struct {
	PrecoProdutos []PrecoProdutos `json:"PrecoProdutos"`
}

func ExtraSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	queryString := strings.ReplaceAll(fmt.Sprintf("https://www.extra.com.br/%s/b?ordenacao=precoCrescente", query), " ", "-")
	queryStringPrice := "https://api.extra.com.br/merchandising/oferta/v1/Preco/Produto/PrecoVenda/?idsProduto="

	if page > 0 {
		queryString += fmt.Sprintf("&page=%d", page)
	}
	if lowPrice != "" && lowPrice != "0" && highPrice != "" && highPrice != "0" {
		queryString += fmt.Sprintf(`&filter=preco^c:%sTO%s`, lowPrice, highPrice)
	} else {
		if lowPrice != "" && lowPrice != "0" {
			queryString += fmt.Sprintf(`&filter=preco^c:%sTO1000000`, lowPrice)
		}
		if highPrice != "" && highPrice != "0" {
			queryString += fmt.Sprintf(`&filter=preco^c:0TO%s`, highPrice)
		}
	}

	html := getExtraRequest(queryString)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Extra [2] Error: ", err)
	}

	htmljson := string(doc.Find(`#__NEXT_DATA__`).First().Text())
	var searchResponse SearchResponse
	err = json.Unmarshal([]byte(htmljson), &searchResponse)
	if err != nil {
		fmt.Println("Extra [3] Error: ", err)
		return productList
	}

	for i, result := range searchResponse.Props.PageProps.InitialState.Search.Results.Products {
		if len(searchResponse.Props.PageProps.InitialState.Search.Results.Products)+1 == i {
			queryStringPrice += result.Id
		} else {
			queryStringPrice += fmt.Sprintf("%s,", result.Id)
		}
	}

	queryStringPrice += fmt.Sprintf("&apiKey=%s", searchResponse.RuntimeConfig.PriceApiKey)

	htmlPrices := getExtraRequest(queryStringPrice)
	var SearchResponsePrices SearchResponsePrices
	err = json.Unmarshal([]byte(htmlPrices), &SearchResponsePrices)
	if err != nil {
		fmt.Println("Extra [4] Error: ", err)
		return productList
	}

	totalPages := 1
	totalProducts := 0
	itensPerPage, _ := strconv.Atoi(searchResponse.RuntimeConfig.ResultsPerPage)
	stringProducts := doc.Find(`p[data-cy="searchCount"]`).First().Text()
	reStringProducts := regexp.MustCompile("[0-9]+")
	matcheStringProduct := reStringProducts.FindAllString(stringProducts, -1)
	if len(matcheStringProduct) > 0 {
		totalProducts, _ = strconv.Atoi(matcheStringProduct[0])
	}
	if (itensPerPage * page) > totalProducts {
		totalPages = page + 1
	}
	for _, result := range searchResponse.Props.PageProps.InitialState.Search.Results.Products {
		productPrice := ""
		for _, priceMap := range SearchResponsePrices.PrecoProdutos {
			priceIDProduct := priceMap.PrecoVenda.IdProduto
			if result.Id == fmt.Sprintf("%d", priceIDProduct) {
				productPrice = strings.TrimSpace(fmt.Sprintf("%f", priceMap.PrecoVenda.Preco))
				break
			}
		}
		if len(productPrice) > 0 {
			productList = append(productList, types.Product{Store: "Extra", TotalPages: fmt.Sprint(totalPages), Image: result.Image, Name: result.Title, Link: result.Href, Stars: fmt.Sprintf("%f", result.Rating), StarsQty: fmt.Sprintf("%f", result.RatingCount), Price: productPrice})
		}
	}
	return productList
}
