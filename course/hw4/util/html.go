package util

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// GetHref comment
func GetHref(t html.Token) string {
	for _, a := range t.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}

	return ""
}

// HTTPGetRequest comment
// func HTTPGetRequest(url string, timeout time.Duration) (*http.Response, context.CancelFunc, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, cancel, err
// 	}

// 	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
// 	if err != nil {
// 		return nil, cancel, err
// 	}

// 	if resp.StatusCode != 200 {
// 		return nil, cancel, fmt.Errorf("status code %d for URL %s", resp.StatusCode, url)
// 	}

// 	return resp, cancel, nil
// }

var textTags = [...]string{
	"address", "article", "aside", "footer", "header",
	"h1", "h2", "h3", "h4", "h5", "h6", "hgroup", "main",
	"nav", "section", "blockquote", "dd", "div", "dl",
	"dt", "hr", "li", "main", "ol", "p", "pre", "ul",
	"a", "abbr", "b", "cite", "dfn", "em", "i", "mark", "q", "s",
	"small", "span", "strong", "sub", "sup", "time", "caption",
	"col", "colgroup", "table", "tbody", "td", "tfoot", "th",
	"thead", "tr",
}

// ExtractTextAndUrls comment
func ExtractTextAndUrls(stream io.Reader) (string, []string, error) {
	var textBuf strings.Builder
	var err error

	docHTML := html.NewTokenizer(stream)
	prevToken := docHTML.Token()
	urlsInHTML := make([]string, 0)

loopEnd:
	for {
		currTokenType := docHTML.Next()

		switch currTokenType {
		case html.ErrorToken:
			err = docHTML.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break loopEnd
			}
		case html.StartTagToken:
			// Check if the token is an <a> tag
			token := docHTML.Token()
			if token.Data == "a" {
				// Extract the href value, if there is one
				if url := GetHref(token); url != "" {
					if strings.HasPrefix(url, "http") {
						urlsInHTML = append(urlsInHTML, url)
					}
				}
			}

			prevToken = token
		case html.TextToken:
			textTagsAsSlice := textTags[:]
			if IsInArray(prevToken.Data, textTagsAsSlice) {
				textInTag := strings.TrimSpace(html.UnescapeString(string(docHTML.Text())))
				if len(textInTag) > 0 {
					_, err = textBuf.WriteString(textInTag)
					if err != nil {
						return "", nil, err
					}
				}
			}
		}
	}

	return textBuf.String(), urlsInHTML, nil
}
