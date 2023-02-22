package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var LC string = "en" // language code
var ExploredPages, PagesInMemory int

// Page is a wikipedia page
type Page struct {
	name     string
	linkWord string
	foundBy  *Page
}

func main() {
	fmt.Println("--- start ---")
	getEnvVars()
	connectDB()

	var startPage, goalPage string
	for {
		fmt.Println("Enter language code (ex: en, pt, es, fr, etc.), see https://en.wikipedia.org/wiki/List_of_Wikipedias")
		fmt.Scan(&LC)
		LC = strings.ToLower(LC)

		fmt.Println("Enter start page name (ex: Penguin)")
		fmt.Scan(&startPage)
		startPage = strings.Replace(strings.ToLower(startPage), " ", "_", -1)

		if len(getPage(startPage)) == 0 {
			fmt.Println("Page does not exist")
			continue
		}

		break
	}

	for {
		fmt.Println("Enter goal page name (ex: Biodiversity)")
		fmt.Scan(&goalPage)
		goalPage = strings.Replace(strings.ToLower(goalPage), " ", "_", -1)

		if startPage == goalPage {
			fmt.Println("Start page and goal page are the same")
			continue
		}

		if len(getPage(goalPage)) == 0 {
			fmt.Println("Page does not exist")
			continue
		}

		break
	}

	// insertPageLinks("Penguin", [][2]string{{"The Mid Lane", "Equator"}, {"album", "Penguin_(album)"}})

	// rows := selectPageLinks("Penguin")

	// for rows.Next() {
	// 	var lc, srcPage, pageName, linkWord string

	// 	// Get values from row.
	// 	err := rows.Scan(&lc, &srcPage, &pageName, &linkWord)
	// 	if err != nil {
	// 		log.Fatal("Error reading rows: ", err.Error())
	// 	}

	// 	fmt.Println(lc, srcPage, pageName, linkWord)
	// }

	// a := getUrl("https://zx.wikipedia.org/wiki/Jadfsi")
	// x := linksInBody(a)
	// fmt.Println(len(x), x)

	list := findWikipediaPath(startPage, goalPage)

	fmt.Printf("Explored pages: %d, Pages in memory: %d", ExploredPages, PagesInMemory)
	fmt.Println(len(list), list)
}

/*
-- input:
lc: language code (ex: en, pt, es, fr, etc.), see https://en.wikipedia.org/wiki/List_of_Wikipedias
start: start page
end: goal page
-- output:
list of wiki pages from start to end [][2]string{pageName, linkWord}
if you click in "linkWord" you will go to "pageName"
*/
func findWikipediaPath(start, end string) [][2]string {
	start = strings.ToLower(start)
	end = strings.ToLower(end)

	// SEARCH (BFS)
	pages := make(map[string]*Page) // map of all pages found so far, key is pageName

	var firstPage, lastPage *Page
	firstPage = &Page{
		name:    start,
		foundBy: nil,
	}

	pages[start] = firstPage

	fifo := []Page{*firstPage}
	for len(fifo) > 0 && lastPage == nil {
		// fifo.pop()
		curPage := fifo[0]
		fifo = fifo[1:]

		pageReferences := getPage(curPage.name)
		for _, pageReference := range pageReferences {
			pageName := strings.ToLower(pageReference[0])
			if pages[pageName] == nil { // this page wasn't found yet
				wordLink := pageReference[1]

				foundPage := Page{
					name:     pageName,
					linkWord: wordLink,
					foundBy:  &curPage,
				}

				if pageName == end {
					lastPage = &foundPage
					break
				}

				pages[pageName] = &foundPage
				fifo = append(fifo, foundPage)
			}
		}

	}

	pagePath := [][2]string{}
	for curPage := lastPage; curPage != nil; curPage = curPage.foundBy {
		pagePath = append(pagePath, [2]string{curPage.name, curPage.linkWord})
	}

	reverseSlice(pagePath)

	return pagePath
}

/*
-- output: [][2]string{pageName, linkWord}
*/
func getPage(pageName string) [][2]string {

	ExploredPages++

	// verify if page is already in database
	rows := selectPageLinks(pageName)
	if rows.Next() { // was found in database

		PagesInMemory++

		links := [][2]string{}
		for rows.Next() {
			var pageName, linkWord string

			// Get values from row.
			err := rows.Scan(&pageName, &linkWord)
			if err != nil {
				log.Fatal("Error reading rows: ", err.Error())
			}

			// fmt.Println(pageName, linkWord)
			links = append(links, [2]string{pageName, linkWord})
		}

		return links
	}

	// Get from wikipedia
	body := getUrl(fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", LC, pageName))
	links := linksInBody(body)

	// insert in database
	insertPageLinks(pageName, links) //TODO: go routine

	return links
}

/*
-- input:
text: the html text of a wikipedia page
-- output:
list of all pages referenced in the text
[2]string{pageName, linkWord}
*/
func linksInBody(text string) [][2]string {
	foundPages := make(map[string]bool) // map of all pages found so far, key is pageName
	pages := [][2]string{}
	urlRegex := regexp.MustCompile(`href="/wiki/([a-zA-Z0-9./?=_-]+)".+?>(.+?)<`)
	for _, page := range urlRegex.FindAllString(text, -1) { // -1 means all matches (doesn't separate the match in paranteses)
		matches := urlRegex.FindStringSubmatch(page) // returns the match in paranteses separated {wholeMatch, pageName, linkWord}

		pageName := strings.ToLower(matches[1])
		linkWord := matches[2]

		if !foundPages[pageName] {
			foundPages[pageName] = true
			pages = append(pages, [2]string{pageName, linkWord})
		}
	}
	return pages
}

/*
-- input:
url: An url
-- output:
the html text of the page
*/
func getUrl(url string) string {
	// a := ""
	// fmt.Print("Will get page: ", url, " (press enter to continue) ")
	// fmt.Scan(&a)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error while getting page: ", url)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error while reading page: ", url)
		log.Fatal(err)
	}

	return string(body)
}

// reverseSlice reverses a slice in place.
func reverseSlice[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// remove duplicates from a slice based on a column
func removeDuplicates(slices [][2]string, column int) [][2]string {
	found := make(map[string]bool)
	result := [][2]string{}

	for _, row := range slices {
		if !found[row[column]] {
			found[row[column]] = true
			result = append(result, row)
		}
	}

	return result
}
