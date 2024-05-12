package builder

import "math/rand"

const randStrLen = 10

func randStr() string {
	var (
		letterRunes  = []rune("abcdefghijklmnopqrstuvwxyz")
		letterLength = len(letterRunes)
		runeString   = make([]rune, randStrLen)
	)

	for i := 0; i < randStrLen; i++ {
		runeString[i] = letterRunes[rand.Intn(letterLength)]
	}

	return string(runeString)
}
