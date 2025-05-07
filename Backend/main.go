package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	fake "price_server/fake"
	Stores "price_server/stores"
	storesinfo "price_server/storesInfo"
	"price_server/types"
	"price_server/utils"
)

type ApiRequest struct {
	Response    chan string
	QueryParams url.Values
}

func handleApiRequest(requestQueue chan *ApiRequest) {
	for {
		apiRequest := <-requestQueue
		response := processRequest(&apiRequest.QueryParams)
		apiRequest.Response <- response
	}
}

func serveHTTP(res http.ResponseWriter, req *http.Request, requestQueue chan *ApiRequest) {
	QueryParams, _ := url.Parse(req.URL.String())

	apiReq := &ApiRequest{
		Response:    make(chan string),
		QueryParams: QueryParams.Query(),
	}

	requestQueue <- apiReq
	resp := <-apiReq.Response
	res.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(res, resp)
}

func main() {
	var path string
	flag.StringVar(&path, "path", "./", "server path")
	flag.Parse()
	maxWorkers := 2
	requestQueue := make(chan *ApiRequest, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go handleApiRequest(requestQueue)
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		serveHTTP(res, req, requestQueue)
	})
	http.HandleFunc("/images/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		http.StripPrefix("/images/", http.FileServer(http.Dir(path+"images/Dark"))).ServeHTTP(res, req)
	})

	http.HandleFunc("/image", func(res http.ResponseWriter, req *http.Request) {
		qParams, _ := url.Parse(req.URL.String())
		resp, err := http.Get(qParams.Query().Get("url"))
		if err != nil {
			fmt.Println(err)
			res.Write([]byte{})
		}
		defer resp.Body.Close()

		imgData, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			res.Write([]byte{})
		}

		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Write(imgData)
	})

	http.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		var storeList []types.Store
		allStores := storesinfo.ListAll()
		for _, storeName := range allStores {
			svgContent := fmt.Sprintf("/images/%s.svg", strings.ToLower(storeName))
			storeList = append(storeList, types.Store{Name: storeName, Logo: svgContent})
		}
		storeJson, err := json.Marshal(storeList)
		if err != nil {
			fmt.Println("Erro ao converter o array:", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(storeJson)
	})

	fmt.Println("Server running on localhost:8000")
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		panic(err)
	}
}

func processRequest(queryParamValue *url.Values) string {
	var fakeRequest = false
	var productList []types.Product
	page := 1
	query := ""
	stores := []string{}
	if len(queryParamValue.Get("q")) != 0 {
		query = queryParamValue.Get("q")
	}
	if len(queryParamValue.Get("pg")) != 0 {
		page, _ = strconv.Atoi(queryParamValue.Get("pg"))
	}
	if len(queryParamValue.Get("stores")) != 0 {
		stores = strings.Split(strings.ToLower(strings.ReplaceAll(queryParamValue.Get("stores"), " ", "")), ",")
	}
	if utils.Contains(stores, "amazon") && storesinfo.GetStoreStatus("amazon") {
		if fakeRequest {
			productList = utils.Concat(productList, fake.AmazonSearchOnPage(query, page))
		} else {
			productList = utils.Concat(productList, Stores.AmazonSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "mercadolivre") && storesinfo.GetStoreStatus("mercadolivre") {
		if fakeRequest {
			productList = utils.Concat(productList, fake.MLSearchOnPage(query, page))
		} else {
			productList = utils.Concat(productList, Stores.MLSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "kabum") && storesinfo.GetStoreStatus("kabum") {
		if fakeRequest {
		} else {
			productList = utils.Concat(productList, Stores.KabumSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "aliexpress") && storesinfo.GetStoreStatus("aliexpress") {
		if fakeRequest {
		} else {
			// not working for now
			productList = utils.Concat(productList, Stores.AliExpressSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "carrefour") && storesinfo.GetStoreStatus("carrefour") {
		if fakeRequest {
		} else {
			// not working for now
			productList = utils.Concat(productList, Stores.CarrefourSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "extra") && storesinfo.GetStoreStatus("extra") {
		if fakeRequest {
		} else {
			productList = utils.Concat(productList, Stores.ExtraSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "casasbahia") && storesinfo.GetStoreStatus("casasbahia") {
		if fakeRequest {
		} else {
			productList = utils.Concat(productList, Stores.CasasBahiaSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "pontofrio") && storesinfo.GetStoreStatus("pontofrio") {
		if fakeRequest {
		} else {
			productList = utils.Concat(productList, Stores.PontoFrioSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if utils.Contains(stores, "magazineluiza") && storesinfo.GetStoreStatus("magazineluiza") {
		if fakeRequest {
		} else {
			productList = utils.Concat(productList, Stores.MagazineLuizaSearchOnPage(query, page, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), strings.ReplaceAll(queryParamValue.Get("pf"), ",", "")))
		}
	}
	if len(queryParamValue.Get("pi")) != 0 {
		pi, err := strconv.ParseFloat(strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""), 64)
		if err != nil {
			fmt.Println(err)
		}
		if pi > 0 {
			productList = utils.PriceHigher(productList, strings.ReplaceAll(queryParamValue.Get("pi"), ",", ""))
		}
	}
	if len(queryParamValue.Get("pf")) != 0 {
		pf, err := strconv.ParseFloat(strings.ReplaceAll(queryParamValue.Get("pf"), ",", ""), 64)
		if err != nil {
			fmt.Println(err)
		}
		if pf > 0 {
			productList = utils.PriceLower(productList, strings.ReplaceAll(queryParamValue.Get("pf"), ",", ""))
		}
	}
	sort.Slice(productList, utils.OrderByPrice(productList))
	if productList == nil {
		productList = []types.Product{}
	}
	productJson, err := json.Marshal(productList)
	if err != nil {
		fmt.Println("Erro ao converter o array:", err)
	}
	return string(productJson)
}
