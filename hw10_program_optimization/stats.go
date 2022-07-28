package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

// easyjson -all stats.go
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)

	result := make(DomainStat)

	reg, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	user := &User{}
	for fileScanner.Scan() {
		err := easyjson.Unmarshal(fileScanner.Bytes(), user)
		if err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}

		if reg.Match([]byte(user.Email)) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
