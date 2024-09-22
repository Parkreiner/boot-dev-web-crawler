package crawler

import (
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "Removes HTTP scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "Removes HTTPS scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "Removes trailing slash",
			inputURL: "blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "Does not transform URL that is already clean",
			inputURL: "blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := normalizeURL(test.inputURL)
			if err != nil {
				t.Errorf(
					"Test %v - '%s' FAIL: unexpected error: %v",
					i,
					test.name,
					err,
				)

				return
			}

			if actual != test.expected {
				t.Errorf(
					"Test %v - %s FAIL: expected URL: %v, actual: %v",
					i,
					test.name,
					test.expected,
					actual,
				)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputUrl  string
		inputBody string
		expected  []string
	}{
		{
			name:     "Translates all relative URLs into absolute URLs",
			inputUrl: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="/path/one">
							<span>Boot.dev</span>
						</a>

						<a href="/path/two">
							<span>Boot.dev</span>
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
		},
		{
			name:     "Doesn't modify absolute URLs",
			inputUrl: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="/path/one">
							<span>Boot.dev</span>
						</a>

						<a href="https://other.com/path/one">
							<span>Boot.dev</span>
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "Returns an empty slice if there are no URLs in the source HTML",
			inputUrl: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<div id="react-spa-app"></div>
					</body>
				</html>
			`,
			expected: []string{},
		},
		{
			name:     "Will automatically remove duplicate URLs present in the HTML",
			inputUrl: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="/path/one">
							<span>Boot.dev</span>
						</a>

						<a href="/path/one">
							<span>Boot.dev</span>
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := getURLsFromHTML(testCase.inputBody, testCase.inputUrl)

			if err != nil {
				t.Errorf(
					"Test %v - '%s' FAIL: unexpected error: %v",
					i,
					testCase.name,
					err,
				)

				return
			}

			if !reflect.DeepEqual(value, testCase.expected) {
				t.Errorf(
					"Test %v - %s FAIL: expected URL: %v, actual: %v",
					i,
					testCase.name,
					testCase.expected,
					value,
				)
			}
		})
	}
}
