package game

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	if IsValid("абспдз") {
		t.Errorf("Ошибка при проверке слова, неверная длина\n")
	}
	if IsValid("asdfg") {
		t.Errorf("Ошибка при проверке слова, неверные символы\n")
	}
	if IsValid("aячсм") {
		t.Errorf("Ошибка при проверке слова, неверный символ\n")
	}
	if IsValid("     ") {
		t.Errorf("Ошибка при проверке слова, неверные символы\n")
	}
	if IsValid("") {
		t.Errorf("Ошибка при проверке слова, неверная длина\n")
	}
	if !IsValid("фывап") {
		t.Errorf("Ошибка при проверке слова, верные символы\n")
	}
	if !IsValid("Ярика") {
		t.Errorf("Ошибка при проверке слова, верные символы\n")
	}
}

func TestHowMuch(t *testing.T) {
	fact := Check("ооооо", "ооооо", CreateDict("ооооо"))

	wait := "ООООО"

	if fact != wait {
		t.Errorf("Ошибка при сравнение слов;\nОжидалось %s;\nДано %s;\n", wait, fact)
	}
	t.Logf("Результат: %s\n", fact)

	fact = Check("топор", "ропот", CreateDict("ропот"))

	wait = "тОПОр"

	if fact != wait {
		t.Errorf("Ошибка при сравнение слов;\nОжидалось %s;\nДано %s;\n", wait, fact)
	}

	t.Logf("Результат: %s\n", fact)

}
