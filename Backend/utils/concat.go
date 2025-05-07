package utils

import "price_server/types"

func Concat(a []types.Product, b []types.Product) []types.Product {
	c := make([]types.Product, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)

	return c
}
