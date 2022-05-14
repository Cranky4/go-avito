package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	runes := []rune(input)
	char := '\n'
	var sb strings.Builder

	for i := 0; i < len(runes); i++ {
		symbol := runes[i]

		number, err := strconv.Atoi(string(symbol))

		switch {
		case number != 0:
			// число пришло раньше буквы
			if char == '\n' {
				return "", ErrInvalidString
			}

			updateResult(&sb, number, char)
			char = '\n'
		case number == 0 && err == nil:
			// поймали 0 в строке - не добавляем текущий символ в результат
			char = '\n'
		default:
			// добавляем букву в результат 1 раз, если к этому моменту не указано число повторений
			if char != '\n' {
				updateResult(&sb, 1, char)
			}

			// проверка на экранирование - следующая цифра выводится, как символ
			if string(symbol) == "\\" {
				if len(runes) < i+1 {
					return "", ErrInvalidString
				}

				// устанавливаем следующий символ, как текущий и сдвигаем проверку символов направо
				char = runes[i+1]
				i++
				continue
			}

			char = symbol
		}
	}

	// обработка последней буквы в строке
	if char != '\n' {
		updateResult(&sb, 1, char)
	}

	return sb.String(), nil
}

// добавить в результат буквы.
func updateResult(sb *strings.Builder, count int, char rune) {
	i := 0
	for i < count {
		sb.WriteRune(char)
		i++
	}
}
