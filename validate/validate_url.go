package validate

import (
	"net/url"
	"strings"

	"github.com/wavly/surf/utils"
)

type InvalidDomainName struct{}
type InvalidDomainFormat struct{}
type DomainTooLong struct{}
type UrlTooLong struct{}
type InvalidUrlPath struct{}
type DomainTooShort struct{}
type InvalidUrlSchema struct{}

func (*InvalidDomainFormat) Error() string {
	return "Only alphabetical characters, digits, and non consecutive hyphens are allowed in the domain name"
}

func (*UrlTooLong) Error() string {
	return "URL is too long, max length is 1000 characters"
}

func (*DomainTooShort) Error() string {
	return "Domain is too short, min length is 4 charecters"
}

func (*InvalidDomainName) Error() string {
	return "URL doesn't contain a valid TLD (Top-Level Domain)"
}

func (*InvalidUrlSchema) Error() string {
	return "Invalid URL schema, only HTTPS schema is allowed"
}

func (*DomainTooLong) Error() string {
	return "Domain Name is too long"
}

func (*InvalidUrlPath) Error() string {
	return "URL path contains invalid characters"
}

func ValidateUrl(link string) (string, error) {
	// Check URL length
	if len(link) > 1000 {
		return link, &UrlTooLong{}
	}

	// Check if the scheme is empty, if so default to https
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	parsedUrl, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	if parsedUrl.Scheme != "https" {
		return link, &InvalidUrlSchema{}
	}

	domain := parsedUrl.Hostname()
	path := parsedUrl.Path

	if err := validateDomain(domain); err != nil {
		return link, err
	}

	// Check for ASCII path characters
	if path != "" {
		if !utils.IsASCII(path) {
			return link, &InvalidUrlPath{}
		}
	}

	return link, nil
}

func validateDomain(domain string) error {
	// Check for TLD
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return &InvalidDomainName{}
	}

	// Check domain length
	if len(domain) > 253 {
		return &DomainTooLong{}
	} else if len(domain) < 4 {
		return &DomainTooShort{}
	}

	// Check for forbidden domain characters & consecutive dashes
	lastDashIndex := 0
	for i, c := range domain {
		if !(utils.IsValidChar(c) || c == '-' || c == '.') || (c == '-' && lastDashIndex+1 == i) {
			return &InvalidDomainFormat{}
		}

		if c == '-' {
			lastDashIndex = i
		}
	}

	if strings.Contains(domain, " ") {
		return &InvalidDomainFormat{}
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
