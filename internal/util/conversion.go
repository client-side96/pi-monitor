package util

import (
	"log"
	"strconv"
)

func ToFloat(stat string) float64 {
	temp, err := strconv.ParseFloat(stat, 64)
	if err != nil {
		log.Fatalf("Cannot convert string to float: %s", err)
	}
	return temp
}
