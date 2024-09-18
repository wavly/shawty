package validate

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/wavly/shawty/utils"
)

type InvalidDomainName struct{}

type DomainTooLong struct{}

type InvalidDomainFormat struct{}

type UrlTooLong struct {
	url uint
}

type InvalidUrlSchema struct {
	schema string
}

func (_ *InvalidDomainFormat) Error() string {
	return "Only alphabetical characters are allowed in the domain name"
}

func (link *UrlTooLong) Error() string {
	return fmt.Sprintf("URL is too long, max lenght is 1000 characters, but got %v", link.url)
}

func (_ *InvalidDomainName) Error() string {
	return "URL doesn't contain a validate TLD (Top-Level Domain)"
}

func (link *InvalidUrlSchema) Error() string {
	return fmt.Sprintf("Invalid URL schema, only HTTPS schema is allowed, but got %s", link.schema)
}

func (_ *DomainTooLong) Error() string {
	return "Domain Name is too long"
}

func ValidateUrl(link string) error {
	if len(link) > 1000 {
		return &UrlTooLong{url: uint(len(link))}
	}

	if !strings.Contains(link, ".") {
		return &InvalidDomainName{}
	}

	split := strings.SplitN(link, ".", 2)
	split[0] = strings.TrimPrefix("https://", split[0])

	if len(split) < 2 || split[0] == "" || split[1] == "" {
		return &InvalidDomainName{}
	}

	if len(split[0]) > 63 || len(split[1]) > 63 {
		return &DomainTooLong{}
	}

	if !utils.IsASCII(split[0]) || !utils.IsASCII(split[1]) {
		return &InvalidDomainFormat{}
	}

	parsedUrl, err := url.Parse(link)
	if err != nil {
		return err
	} else if parsedUrl.Scheme != "https" {
		return &InvalidUrlSchema{schema: parsedUrl.Scheme}
	}

	return nil
}
