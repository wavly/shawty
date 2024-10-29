package validate

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/wavly/shawty/utils"
)

type InvalidDomainName struct{}
type InvalidDomainFormat struct{}
type DomainTooLong struct{}
type DomainTooShort struct {
	domain string
}

type UrlTooLong struct {
	url uint
}
type InvalidUrlSchema struct {
	schema string
}
type InvalidUrlPath struct {
	path string
}

func (*InvalidDomainFormat) Error() string {
	return "Only alphabetical characters, digits, and non consecutive hyphens are allowed in the domain name"
}

func (link *UrlTooLong) Error() string {
	return fmt.Sprintf("URL is too long, max length is 1000 characters, but got %v", link.url)
}

func (link *DomainTooShort) Error() string {
	return fmt.Sprintf("Domain is too short, min length is 4 charecters, domain: %s", link.domain)
}

func (*InvalidDomainName) Error() string {
	return "URL doesn't contain a valid TLD (Top-Level Domain)"
}

func (link *InvalidUrlSchema) Error() string {
	return fmt.Sprintf("Invalid URL schema, only HTTPS schema is allowed, but got %s", link.schema)
}

func (*DomainTooLong) Error() string {
	return "Domain Name is too long"
}

func (link *InvalidUrlPath) Error() string {
	return fmt.Sprintf("URL path contains invalid characters: %s", link.path)
}

func ValidateUrl(link string) (string, error) {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	// Check if the scheme is empty, if so default to https
	if parsedUrl.Scheme == "" {
		link = "https://" + link
		parsedUrl, err = url.Parse(link)
		if err != nil {
			return link, err
		}
	}

	if parsedUrl.Scheme != "https" {
		return link, &InvalidUrlSchema{schema: parsedUrl.Scheme}
	}

	// Check URL length
	if len(link) > 1000 {
		return link, &UrlTooLong{url: uint(len(link))}
	}

	domain := parsedUrl.Hostname()
	path := parsedUrl.Path

	if err := validateDomain(domain); err != nil {
		return link, err
	}

	// Check for ASCII path characters
	if path != "" {
		if !utils.IsASCII(path) {
			return link, &InvalidUrlPath{path: path}
		}
	}

	return link, nil
}

func validateDomain(domain string) error {
	// Check domain length
	if len(domain) > 253 {
		return &DomainTooLong{}
	} else if len(domain) < 4 {
		return &DomainTooShort{domain: domain}
	}

	// Check for allowed domain characters
	for _, c := range domain {
		if !(utils.IsValidChar(c) || c == '-' || c == '.') {
			return &InvalidDomainFormat{}

		}
	}
	if strings.Contains(domain, " ") {
		return &InvalidDomainFormat{}
	}

	// Check for consecutive dashes
	re := regexp.MustCompile(`-{2,}`)
	if re.MatchString(domain) {
		return &InvalidDomainFormat{}
	}

	// Check for TLD
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return &InvalidDomainName{}
	}

	// Validate each domain part seperately
	for _, part := range domainParts {
		if err := isValidDomainPart(part); err != nil {
			return err
		}
	}

	return nil
}

func isValidDomainPart(part string) error {
	if len(part) > 63 {
		return &DomainTooLong{}
	}

	// Check for leading or trailing dashes & empty parts
	if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") || part == "" {
		return &InvalidDomainFormat{}
	}

	return nil
}
