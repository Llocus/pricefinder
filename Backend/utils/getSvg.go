package utils

import "io/ioutil"

func GetSVG(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}
