package game

import "log"

func IsValid(answer string) bool {
	word := []rune(answer)
	if len(word) != 5 {
		return false
	}
	for i := 0; i < 5; i++ {
		if word[i] > 1103 || word[i] < 1040 {
			return false
		}
	}
	return true
}

func CreateDict(right string) map[rune]int {
	correct := []rune(right)
	res := make(map[rune]int)
	for i := 0; i < 5; i++ {
		if _, ok := res[Low(correct[i])]; !ok {
			res[Low(correct[i])] = 0
		}
		res[Low(correct[i])]++
	}
	return res
}

func Low(r rune) rune {
	if r < 1072 {
		return r + 32
	}
	return r
}

func High(r rune) rune {
	if r >= 1072 {
		return r - 32
	}
	return r
}

func Check(answer, right string, d map[rune]int) string {
	word := []rune(answer)
	correct := []rune(right)

	res := make([]rune, 5)

	for i := 0; i < 5; i++ {
		if Low(word[i]) == Low(correct[i]) {
			d[Low(word[i])]--
			res[i] = High(word[i])
			log.Println(High(word[i]))
		}
	}
	for i := 0; i < 5; i++ {
		if d[Low(word[i])] != 0 && Low(word[i]) != Low(correct[i]) {
			d[Low(word[i])]--
			res[i] = Low(word[i])
		} else if Low(word[i]) != Low(correct[i]) {
			res[i] = []rune("_")[0]
		}
	}
	return string(res)
}

// маленькие буквы русского алфавита: 1072-1103
// Большие буквы русского алфавита: 1040-1071

func Success(s string) bool {
	r := []rune(s)
	for i := 0; i < 5; i++ {
		if r[i] == rune('_') || r[i] == Low(r[i]) {
			return false
		}
	}
	return true
}
