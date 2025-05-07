package storesinfo

import "strings"

func GetStoreStatus(name string) bool {
	if strings.ToLower(name) == "amazon" {
		return true
	}
	if strings.ToLower(name) == "mercadolivre" {
		return true
	}
	if strings.ToLower(name) == "aliexpress" {
		return false
	}
	if strings.ToLower(name) == "shopee" {
		return false
	}
	if strings.ToLower(name) == "shein" {
		return false
	}
	if strings.ToLower(name) == "kabum" {
		return true
	}
	if strings.ToLower(name) == "casasbahia" {
		return true
	}
	if strings.ToLower(name) == "magazineluiza" {
		return true
	}
	if strings.ToLower(name) == "pontofrio" {
		return true
	}
	if strings.ToLower(name) == "americanas" {
		return false
	}
	if strings.ToLower(name) == "extra" {
		return true
	}
	if strings.ToLower(name) == "girafa" {
		return false
	}
	if strings.ToLower(name) == "leroymerlin" {
		return false
	}
	if strings.ToLower(name) == "philco" {
		return false
	}
	if strings.ToLower(name) == "consul" {
		return false
	}
	if strings.ToLower(name) == "netshoes" {
		return false
	}
	if strings.ToLower(name) == "fastshop" {
		return false
	}
	if strings.ToLower(name) == "nike" {
		return false
	}
	if strings.ToLower(name) == "carrefour" {
		return false
	}
	if strings.ToLower(name) == "centauro" {
		return false
	}
	if strings.ToLower(name) == "polishop" {
		return false
	}
	if strings.ToLower(name) == "madeiramadeira" {
		return false
	}
	if strings.ToLower(name) == "olympikus" {
		return false
	}
	return false
}
