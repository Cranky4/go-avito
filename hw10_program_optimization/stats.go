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
	// ID       int
	// Name     string
	// Username string
	Email string
	// Phone    string
	// Password string
	// Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)
	var i int
	for fileScanner.Scan() {
		user := &User{}
		if err = easyjson.Unmarshal([]byte(fileScanner.Text()), user); err != nil {
			return
		}
		result[i] = *user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	r, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		if r.Match([]byte(user.Email)) {
			d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[d]++
		}
	}

	return result, nil
}
