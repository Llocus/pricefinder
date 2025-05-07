package Stores

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"price_server/types"
	"price_server/utils"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getAmazonTotalPages(html string) int {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Amazon [0] Error: ", err)
	}
	var pageArr []string
	page := ""
	doc.Find(".s-pagination-strip").Children().Each(func(i int, s *goquery.Selection) {
		s.Parent().Find(".s-pagination-item").Each(func(j int, s2 *goquery.Selection) {
			if j == (s.Parent().Find(".s-pagination-item").Length() - 2) {
				pageArr = append(pageArr, strings.Replace(s2.Text(), " ", "", -1))
			}
		})
	})
	if len(pageArr) > 0 {
		page = pageArr[0]
	}
	pageNum, _ := strconv.Atoi(page)
	return pageNum
}

func extractAmazonSearch(html string) []types.Product {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Amazon [1] Error: ", err)
	}
	var data []types.Product
	doc.Find(".s-result-list .puis-card-container").Each(func(i int, s *goquery.Selection) {
		nameF := s.Find(".a-link-normal span.a-size-base-plus")
		name, _ := nameF.Html()
		linkF, _ := s.Find(".s-product-image-container").Html()
		link := ""
		stars := ""
		starsQty := ""
		reviewSection, _ := s.Html()
		section, _ := s.Html()

		re1 := regexp.MustCompile(`(?i)title="([^"]*)"`)
		matches1 := re1.FindStringSubmatch(name)
		re2 := regexp.MustCompile(`(?i)href="([^"]*)"`)
		matches2 := re2.FindStringSubmatch(linkF)
		re3 := regexp.MustCompile(`(?i)class="a-icon-alt">([^"]*) de 5 estrelas`)
		matches3 := re3.FindStringSubmatch(section)
		re4 := regexp.MustCompile(`(?i)<span class="a-size-base s-underline-text">([^"]*)</span> </a>`)
		matches4 := re4.FindStringSubmatch(reviewSection)

		if len(matches1) > 1 {
			name = matches1[1]
		}
		if len(matches2) > 1 {
			link = utils.DecodeURL(matches2[1])
		}
		if len(matches3) > 1 {
			stars = strings.Replace(matches3[1], ",", ".", -1)
		}
		if len(matches4) > 1 {
			starsQty = matches4[1]
		}
		imageF := s.Find(".s-product-image-container")
		image, _ := imageF.Html()
		re := regexp.MustCompile(`(?i)src="([^"]*)"`)
		matches := re.FindStringSubmatch(image)
		if len(matches) > 1 {
			image = matches[1]
		}

		price := strings.ReplaceAll(s.Find(".a-link-normal > .a-price > .a-offscreen").Text(), " ", "")
		fakeLink := false

		reLink := regexp.MustCompile(`sspa/click`)
		linkMatches := reLink.FindStringSubmatch(link)
		if len(linkMatches) > 1 {
			fakeLink = true
		}

		starsNum, _ := strconv.ParseFloat(stars, 64)
		if starsNum == 0 {
			stars = "0"
		}
		if starsNum > 3 && price != "" && !fakeLink {
			priceArr := strings.Split(price, "R$ ")
			if len(priceArr) >= 1 {
				price = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(priceArr[0], "R$", ""), " ", ""), ".", ""), ",", "."))
				price = strings.TrimSpace(strings.Split(strings.ReplaceAll(price, "\u00a0", " "), " ")[0])
			}
			data = append(data, types.Product{Store: "Amazon", TotalPages: fmt.Sprintf("%d", getAmazonTotalPages(html)), Image: image, Name: name, Link: "https://amazon.com.br" + link, Stars: stars, StarsQty: starsQty, Price: price})
		}
	})
	sort.Slice(data, utils.OrderByPrice(data))
	return data
}

func getAmazonRequest(url string) string {
	data := []byte(``)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Amazon [2] Error: ", err)
		fmt.Println("Erro ao criar a requisição:", err)
		return ""
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Amazon [3] Error: ", err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)
	}
	return ""
}

func AmazonSearchOnPage(query string, page int, lowPrice string, highPrice string) []types.Product {
	var productList []types.Product
	queryString := fmt.Sprintf("https://www.amazon.com.br/s?k=%s&page=%d&s=price-asc-rank", strings.Replace(query, " ", "+", -1), page)
	queryString += "&sr_nr_p_72_3"
	if lowPrice != "" && lowPrice != "0" {
		queryString += fmt.Sprintf("&low-price=%s", lowPrice)
	}
	if highPrice != "" && highPrice != "0" {
		queryString += fmt.Sprintf("&high-price=%s", highPrice)
	}
	html := getAmazonRequest(queryString)
	productList = extractAmazonSearch(html)
	return productList
}
