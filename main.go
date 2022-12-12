package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Book stores information
type Book struct {
	Title          string
	URL            string
	Author         string
	ISBN           string
	Price          int64
	TotalPage      int64
	CollectionType string
	Description    string
	Rating         float32
}

func main() {
	fName := "books.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: nxbkimdong.com.vn, www.nxbkimdong.com.vn
		colly.AllowedDomains("nxbkimdong.com.vn", "www.nxbkimdong.com.vn"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./nxbkimdong_book_cache"),
	)

	// Create another collector to scrape book details
	detailCollector := c.Clone()

	books := make([]Book, 0, 200)

	// On every <a> element which has "href" attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		// If link is not all page then return
		if !strings.Contains(link, "collections/all?page=") {
			return
		}

		// start scaping the page under the link found
		e.Request.Visit(link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every <a> element with collection-product-card class call callback
	c.OnHTML(`.product-item > .product-img > a`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains "nxbkimdong.com.vn/products"
		bookDetailUrl := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(bookDetailUrl, "nxbkimdong.com.vn/products") {
			detailCollector.Visit(bookDetailUrl)
		}
	})

	// Extract details of the book
	detailCollector.OnHTML(`section[id=product-wrapper]`, func(e *colly.HTMLElement) {
		log.Println("Book found", e.Request.URL)
		title := e.ChildText("div.header_wishlist > h1")
		if title == "" {
			log.Println("No title found", e.Request.URL)
			return
		}

		price, err := strconv.ParseInt(ClearString(e.ChildText(".ProductPrice")), 10, 0)
		if err != nil {
			price = -1
		}

		pageNumber, err := strconv.ParseInt(ClearString(e.ChildText("ul>li:nth-child(5)")), 10, 0)
		if err != nil {
			pageNumber = 0
		}

		book := Book{
			Title:          title,
			URL:            e.Request.URL.String(),
			Author:         e.ChildText("ul>li:nth-child(2) > a"),
			ISBN:           e.ChildText("ul>li:nth-child(1) > strong"),
			Price:          price,
			TotalPage:      pageNumber,
			CollectionType: e.ChildText("ul>li:nth-child(8) > a"),
			Description:    "",
			Rating:         0,
		}

		books = append(books, book)
	})

	c.Visit("https://nxbkimdong.com.vn/collections/all?page=1")

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(books)

	log.Println("Finished")
}

var nonAlphanumericRegex = regexp.MustCompile(`[^0-9]+`)

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
