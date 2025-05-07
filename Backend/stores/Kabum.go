package Stores

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"price_server/types"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func getKabumRequest(url string) string {
	data := []byte(``)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Kabum [0] Error: ", err)
		return ""
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Kabum [1] Error: ", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		return string(body)
	}
	return ""
}

func removeAccentsAndReplaceSpaces(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	res, _, _ := transform.String(t, s)
	s = res
	s = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(s, "")
	return "/" + strings.ToLower(strings.ReplaceAll(string(s), " ", "-"))
}

type KabumMeta struct {
	TotalItemsCount int `json:"total_items_count"`
	TotalPagesCount int `json:"total_pages_count"`
}

type KabumData struct {
	Links      KabumLinks      `json:"links"`
	Attributes KabumAttributes `json:"attributes"`
}

type KabumAttributes struct {
	Title              string       `json:"title"`
	Price              float64      `json:"price"`
	Rating             float64      `json:"score_of_ratings"`
	NumRating          float64      `json:"number_of_ratings"`
	DiscountPercentage float64      `json:"discount_percentage"`
	PriceWithDiscount  float64      `json:"price_with_discount"`
	Available          bool         `json:"available"`
	Images             []string     `json:"images"`
	Photos             KabumPhotosP `json:"photos"`
}

type KabumPhotosP struct {
	P  []string `json:"p"`
	M  []string `json:"m"`
	G  []string `json:"g"`
	GG []string `json:"gg"`
}

type KabumLinks struct {
	Self string `json:"self"`
}

type KabumSearchResponse struct {
	Meta KabumMeta   `json:"meta"`
	Data []KabumData `json:"data"`
}

func StringToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func KabumSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	queryString := fmt.Sprintf("https://servicespub.prod.api.aws.grupokabum.com.br/catalog/v2/search?query=%s&page_size=20&sort=price", strings.Replace(query, " ", "-", -1))
	if page > 0 {
		queryString += fmt.Sprintf("&page_number=%d", page)
	}
	if lowPrice != "" && lowPrice != "0" && highPrice != "" && highPrice != "0" {
		queryString += "&facet_filters=" + StringToBase64(fmt.Sprintf(`{"kabum_product":["true"],"price":{"min":%s,"max":%s}}`, lowPrice, highPrice))
	} else {
		if lowPrice != "" && lowPrice != "0" {
			queryString += "&facet_filters=" + StringToBase64(fmt.Sprintf(`{"kabum_product":["true"],"price":{"min":%s,"max":1000000}}`, lowPrice))
		}
		if highPrice != "" && highPrice != "0" {
			queryString += "&facet_filters=" + StringToBase64(fmt.Sprintf(`{"kabum_product":["true"],"price":{"min":0,"max":%s}}`, highPrice))
		}
	}
	html := getKabumRequest(queryString)

	var searchResponse KabumSearchResponse
	err := json.Unmarshal([]byte(html), &searchResponse)
	if err != nil {
		fmt.Println("Kabum [2] Error: ", err)
		return productList
	}

	for _, result := range searchResponse.Data {
		formatLink := "https://www.kabum.com.br" + strings.Replace(strings.ReplaceAll(result.Links.Self, "/catalog/v2", ""), "products", "produto", 1) + removeAccentsAndReplaceSpaces(result.Attributes.Title)
		stars := fmt.Sprintf("%d", 0)
		numStars := fmt.Sprintf("%d", 0)
		if result.Attributes.Rating > 0 {
			stars = fmt.Sprintf("%.1f", result.Attributes.Rating)
			numStars = fmt.Sprintf("%f", result.Attributes.NumRating)
		}
		if result.Attributes.Rating >= 3 {
			productList = append(productList, types.Product{Store: "Kabum", TotalPages: fmt.Sprintf("%d", searchResponse.Meta.TotalPagesCount), Image: result.Attributes.Images[0], Name: result.Attributes.Title, Link: formatLink, Stars: stars, StarsQty: numStars, Price: fmt.Sprintf("%f", result.Attributes.Price)})
		}
	}
	return productList
}
