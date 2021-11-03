package fastly

import "math/rand"

const randNameSeed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randNameSeed[rand.Intn(len(randNameSeed))]
	}
	return string(b)
}
