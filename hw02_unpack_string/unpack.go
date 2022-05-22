package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	const empty = '\n'

	runes := []rune(input)
	char := empty
	var sb strings.Builder

	for i := 0; i < len(runes); i++ {
		symbol := runes[i]

		number, err := strconv.Atoi(string(symbol))

		switch {
		case number != 0:
			// число пришло раньше буквы
			if char == empty {
				return "", ErrInvalidString
			}

			updateResult(&sb, number, char)
			char = empty
		case number == 0 && err == nil:
			// число пришло раньше буквы
			if char == empty {
				return "", ErrInvalidString
			}
			// поймали 0 в строке - не добавляем текущий символ в результат
			char = empty
		default:
			// добавляем букву в результат 1 раз, если к этому моменту не указано число повторений
			if char != empty {
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
	if char != empty {
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
