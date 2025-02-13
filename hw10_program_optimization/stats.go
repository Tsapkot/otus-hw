package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigFastest

	for scanner.Scan() {
		line := scanner.Bytes()

		var user struct {
			Email string `json:"email"`
		}
		if err := json.Unmarshal(line, &user); err != nil || user.Email == "" {
			continue
		}

		emailParts := strings.SplitN(user.Email, "@", 2)
		if len(emailParts) < 2 {
			continue
		}

		domainPart := strings.ToLower(emailParts[1])

		if strings.HasSuffix(domainPart, "."+domain) {
			result[domainPart]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading input: %w", err)
	}

	return result, nil
}
