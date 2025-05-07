package storesinfo

import "strings"

func ListAll() []string {
	list := "Amazon,MercadoLivre,Aliexpress,Shopee,Shein,Kabum,CasasBahia,MagazineLuiza,PontoFrio,Americanas,Girafa,LeroyMerlin,Philco,Extra,Consul,Netshoes,FastShop,Nike,Carrefour,Centauro,Polishop,MadeiraMadeira,Olympikus"
	var result []string
	for _, store := range strings.Split(list, ",") {
		if GetStoreStatus(store) {
			result = append(result, store)
		}
	}
	return result
}
