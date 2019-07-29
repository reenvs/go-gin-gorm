package util

import "math/rand"

func Random(ids []uint32) {
	if len(ids) <= 0 {
		return
	}
	for i := len(ids) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		ids[i], ids[num] = ids[num], ids[i]
	}
}
