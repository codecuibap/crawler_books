package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Book stores information
type Book struct {
	Title          string  `json:"title"`
	URL            string  `json:"detail"`
	Author         string  `json:"author"`
	ISBN           string  `json:"isbn"`
	Price          int64   `json:"price"`
	TotalPage      int64   `json:"number_of_page"`
	CollectionType string  `json:"category"`
	Name           string  `json:"book_nane"`
	Group          string  `json:"book_shelf"`
	Description    string  `json:"description"`
	Rating         float32 `json:"rate"`
}

func main() {
	site := flag.String("site", "nxbkimdong.com.vn", "input site to scrap")
	flag.Parse()

	config, err := LoadConfig(*site)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fName := fmt.Sprintf("%s_books.json", *site)
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: nxbkimdong.com.vn, www.nxbkimdong.com.vn
		colly.AllowedDomains(config.Site, fmt.Sprintf("www%s", config.Site)),

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
		if !strings.Contains(link, config.Collection) {
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
	c.OnHTML(config.Product, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains "nxbkimdong.com.vn/products"
		bookDetailUrl := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(bookDetailUrl, config.UrlDetail) {
			log.Println("Load child: ", bookDetailUrl)
			detailCollector.Visit(bookDetailUrl)
		}
	})

	// Extract details of the book
	detailCollector.OnHTML(config.Section, func(e *colly.HTMLElement) {
		log.Println("Book found", e.Request.URL)
		title := e.ChildText(config.Title)
		if title == "" {
			log.Println("No title found", e.Request.URL)
			return
		}

		price, err := strconv.ParseInt(ClearString(e.ChildText(config.Price)), 10, 0)
		if err != nil {
			price = -1
		}

		pageNumber, err := strconv.ParseInt(ClearString(e.ChildText(config.Page)), 10, 0)
		if err != nil {
			pageNumber = 0
		}

		book := Book{
			Title:          title,
			URL:            e.Request.URL.String(),
			Author:         e.ChildText(config.Author),
			ISBN:           e.ChildText(config.ISBN),
			Price:          price,
			TotalPage:      pageNumber,
			CollectionType: e.ChildText(config.Category),
			Name:           e.ChildText(config.Name),
			Group:          e.ChildText(config.Group),
			Description:    e.ChildText(config.Description),
			Rating:         0,
		}
		log.Println("Book", book)
		books = append(books, book)
	})

	c.Visit(config.ScrapSite)

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output

	enc.Encode(books)

	log.Println("Finished", len(books))
}

var nonAlphanumericRegex = regexp.MustCompile(`[^0-9]+`)

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
