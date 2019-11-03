package common

import "math/rand"

func GetRand(num int) int {
	return rand.Intn(num)
}
