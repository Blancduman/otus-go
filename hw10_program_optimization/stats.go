package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	stats := make(DomainStat)

	for scanner.Scan() {
		var user User

		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}

		host := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		if host == "" {
			continue
		}

		dmn := strings.ToLower(strings.SplitN(host, ".", 2)[1])
		if dmn == strings.ToLower(domain) {
			stats[host]++
		}
	}

	return stats, nil
}
