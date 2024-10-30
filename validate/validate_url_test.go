package validate

import (
	"testing"
)

func TestValidateURLs(t *testing.T) {
	tests := []struct {
		link string
		err  error
	}{
		{"https://google.com", nil},
		{"https://www.example.com", nil},
		{"https://sub1.sub2.example.com", nil},
		{"https://example.com/path", nil},
		{"https://example.com/path/to/resource", nil},
		{"https://example.com/path?param=value", nil},
		{"https://example.com:443", nil},
		{"https://example.com/path#fragment", nil},
		{"https://example.com/path-with-dash", nil},
		{"example.com", nil},
		{"sub.example.com", nil},
		{"example.com/path", nil},
		{"test-dash-domain.com", nil},
		{"microsoft.com", nil},
		{"https://user:pass@ple.com", nil},
		{"github.com", nil},
		{"stackoverflow.com", nil},
		{"medium.com", nil},
		{"amazon.com", nil},
		{"netflix.com", nil},

		// Invalid Schemes
		{"http://example.com", &InvalidUrlSchema{}},
		{"ftp://example.com", &InvalidUrlSchema{}},
		{"ws://example.com", &InvalidUrlSchema{}},
		{"file://example.com", &InvalidUrlSchema{}},
		{"mailto:user@example.com", &InvalidUrlSchema{}},
		{"fiffdsale://something.com", &InvalidUrlSchema{}},

		// Invalid Domain Format
		{"https://example..com", &InvalidDomainFormat{}},

		{"xn--bcher-kva.com", &InvalidDomainFormat{}},
		{"to-o--many---dashes.com", &InvalidDomainFormat{}},

		{"https://example--com", &InvalidDomainName{}},
		{"https://example-.com", &InvalidDomainFormat{}},
		{"https://-example.com", &InvalidDomainFormat{}},
		{"https://exam!ple.com", &InvalidDomainFormat{}},
		{"https://exam#ple.com", &InvalidDomainName{}},
		{"https://exam$ple.com", &InvalidDomainFormat{}},
		{"https://exam&ple.com", &InvalidDomainFormat{}},
		{"https://a..c", &InvalidDomainFormat{}},

		// Domain Too Short
		{"https://x.y", &DomainTooShort{}},
		{"https://a.c", &DomainTooShort{}},
		{"https://xyz", &InvalidDomainName{}},

		// Invalid Domain Name (Missing TLD)
		{"https://localhost", &InvalidDomainName{}},
		{"https://internal", &InvalidDomainName{}},
		{"https://example", &InvalidDomainName{}},
		{"https://com", &InvalidDomainName{}},

		// Invalid Characters in Path
		{"https://example.com/path/√∆˚", &InvalidUrlPath{}},
		{"https://example.com/path/漢字", &InvalidUrlPath{}},
		{"https://example.com/path/🚀", &InvalidUrlPath{}},
		{"https://example.com/path/ñ", &InvalidUrlPath{}},

		// Edge Cases
		{"https://", &InvalidDomainName{}},
		{"https:///path", &InvalidDomainName{}},
		{"https://.", &DomainTooShort{}},
		{"https://.com", &InvalidDomainFormat{}},
		{"", &InvalidDomainName{}},
		{"https://example.com/ path", nil}, // Space in path is actually valid when encoded
		{"https://example.com/path/", nil},
		{"https://example.com?param=value", nil},
		{"https://example.com#fragment", nil},

		// More Valid URLs with Various TLDs
		{"https://example.co.uk", nil},
		{"https://example.com.br", nil},
		{"https://example.io", nil},
		{"https://example.dev", nil},
		{"https://example.app", nil},
		{"https://example.cloud", nil},

		// More Invalid Cases
		{"https://.example.com", &InvalidDomainFormat{}},
		{"https://example.com.", &InvalidDomainFormat{}},
		{"https://exa....mple.com", &InvalidDomainFormat{}},

		// Additional Edge Cases
		{"https://123.123.123.123", nil},
		{"https://example.museum", nil},
		{"https://example.travel", nil},
		{"https://example.co.uk.com", nil},
		{"https://test.example.co.uk", nil},
	}

	for _, tt := range tests {
		_, err := ValidateUrl(tt.link)
		if err != tt.err {
			t.Errorf("ValidateUrl(%s): want: %v | %v", tt.link, tt.err, err)
		}
	}
}
