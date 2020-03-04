package helpers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func GetDomain(link string) (string, error) {
	if !(strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")) {
		link = "http://" + link
	}

	u, err := ParseURL(link)
	if err != nil {
		return "", err
	}

	return u.Hostname(), nil
}

//
// Returns unique items in a slice
//
func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func Exists(name string) bool {
	f, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	if f.Size() <= 1 && !f.IsDir() {
		return false
	}
	return true
}

func ParseURL(path string) (*url.URL, error) {
	URL, err := url.Parse(strings.TrimSpace(path))
	if err != nil {
		return nil, err
	}

	return URL, nil
}

func Bfs(nav *goquery.Selection, crc map[string]int, action func(node *html.Node, crc map[string]int)) {
	queue := make([]*goquery.Selection, 0)
	queue = append(queue, nav)

	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]

		action(nextUp.Nodes[0], crc)

		nextUp.Children().Each(func(_ int, selection *goquery.Selection) {
			queue = append(queue, selection)
		})
	}
}

func ReadPathHTML(filepath string) (*goquery.Document, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(b)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	return doc, err
}

func CheckNodeData(node html.Node) bool {
	if node.Data == "" ||
		node.Data == "script" ||
		node.Data == "noscript" ||
		node.Data == "hr" ||
		node.Data == "select" ||
		node.Data == "option" ||
		node.Data == "style" ||
		node.Data == "meta" ||
		node.Data == "link" ||
		node.Data == "head" {
		return false
	}

	return true
}

// OutputHTML returns the text including tags name.
func OutputHTML(node *html.Node, self bool) string {
	var buf bytes.Buffer
	if self {
		_ = html.Render(&buf, node)
	} else {
		for n := node.FirstChild; n != nil; n = n.NextSibling {
			_ = html.Render(&buf, n)
		}
	}
	return buf.String()
}

func StringMinifier(in string) (out string) {
	white := false
	in = strings.TrimSpace(in)
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out += out
			}
			white = true
		} else {
			out += string(c)
			white = false
		}
	}
	return
}

func CreateFile(filepath, data string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Sync()
	}()

	return nil
}
