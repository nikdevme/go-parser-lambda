package parsers

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-parser-lambda/pkg/aws"
	"go-parser-lambda/pkg/shared"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ParseProductsList(client *http.Client) ([]shared.Product, string) {
	pages := [3]string{"PageOne", "PageTwo", "PageThree"}
	var products []shared.Product

	// HashTable for uniq item
	hashTable := make(map[string]string)

	// HashTable for max date
	hashTableDate := make(map[string]int)

	for _, page := range pages {
		// set parameters
		req, err := http.NewRequest("GET", "https://examples.com/pages/"+page, nil)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Set parametr for request page error: %v", err))
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Request
		resp, err := client.Do(req)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Request page error: %v", err))
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				aws.UploadLogToS3("dataOne", fmt.Sprintf("io.ReadCloser close error: %v", err))
			}
		}(resp.Body)

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Read response page error: %v", err))
		}

		// Parse response to string and to goquery
		reader := strings.NewReader(string(body))
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Parse response to string and to goquery error: %v", err))
			log.Fatal(err)
		}
		doc.Find("tr.Item").Each(func(i int, s *goquery.Selection) {
			var product shared.Product

			date, _ := s.Attr("date")
			name, _ := s.Attr("Item")
			parameter1, _ := s.Attr("parameter1")
			parameter2 := s.Find("parameter2").Text()
			parameter3 := s.Find("parameter3").Find("span").First().Text()

			// change time format
			layout := "1/2/2006 3:04:05 PM"
			dateTime, err := time.Parse(layout, date)
			if err != nil {
				aws.UploadLogToS3("dataOne", fmt.Sprintf("Date parsing error: %v", err))
				return
			}
			unixTimeMs := dateTime.UnixNano() / int64(time.Millisecond)
			strDate := dateTime.Format("02-01-2006")

			// set markets
			product.Page = page
			product.Date = strconv.FormatInt(unixTimeMs, 10)
			product.Name = name
			product.Parameter1 = parameter1
			product.Parameter2 = parameter2
			product.Parameter3 = parameter3

			// Check uniq item
			if _, ok := hashTable[name]; !ok {
				hashTable[name] = ""
				products = append(products, product)
			}

			// Counter dates
			if _, ok := hashTableDate[strDate]; ok {
				hashTableDate[strDate] = hashTableDate[strDate] + 1
			} else {
				hashTableDate[strDate] = 1
			}

		})

	}

	// Find max date
	var maxDate string
	maxCount := 0

	for date, count := range hashTableDate {
		if count > maxCount {
			maxCount = count
			maxDate = date
		}
	}

	return products, maxDate
}

func ParseDescription(client *http.Client, products []shared.Product, parseURL string) []shared.ProductDescription {
	var prodDesc []shared.ProductDescription

	for _, product := range products {
		var productIn shared.ProductArray
		var productOut shared.ProductDescription

		// set parameters
		url := fmt.Sprintf("%sName=%s&Date=%s", parseURL, product.Name, product.Date)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Set parameter for request page description error: %v", err))
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Request
		resp, err := client.Do(req)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Request description page error: %v", err))
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				aws.UploadLogToS3("dataOne", fmt.Sprintf("io.ReadCloser close error: %v", err))
			}
		}(resp.Body)

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Read response description page error: %v", err))
		}

		// Parse response to string and to goquery
		reader := strings.NewReader(string(body))
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("Parse description response to string and to goquery error: %v", err))
			log.Fatal(err)
		}

		//fmt.Printf(doc)
		err = json.Unmarshal([]byte(doc.Text()), &productIn)

		productOut.Page = product.Page
		productOut.Date = product.Date
		productOut.Name = product.Name
		productOut.Parameter1 = product.Parameter1
		productOut.Parameter2 = product.Parameter2
		productOut.Parameter3 = product.Parameter3

		prodDesc = append(prodDesc, productOut)
		time.Sleep(300 * time.Millisecond)

	}

	return prodDesc
}

func ParseDataOneProduct(client *http.Client, market shared.Product, parseURL string) string {

	// set parameters 1
	url := fmt.Sprintf("%s?Name=%s&Date=%s", parseURL, market.Name, market.Date)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		aws.UploadLogToS3("dataOne", fmt.Sprintf("Set parametr for request error: %v", err))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Request
	resp, err := client.Do(req)
	if err != nil {
		aws.UploadLogToS3("dataOne", fmt.Sprintf("Request one product page error: %v", err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			aws.UploadLogToS3("dataOne", fmt.Sprintf("io.ReadCloser close error: %v", err))
		}
	}(resp.Body)

	// Read response
	body1, err := io.ReadAll(resp.Body)
	if err != nil {
		aws.UploadLogToS3("dataOne", fmt.Sprintf("Read response one product page error: %v", err))
	}
	time.Sleep(300 * time.Millisecond)

	return string(body1)
}
