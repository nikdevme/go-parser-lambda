package main

import (
	"encoding/json"
	"fmt"
	"go-parser-lambda/pkg/auth"
	"go-parser-lambda/pkg/aws"
	"go-parser-lambda/pkg/files"
	"go-parser-lambda/pkg/parsers"
	"go-parser-lambda/pkg/shared"
	"go-parser-lambda/pkg/utils"
	"os"
)

func main() {
	authURL := os.Getenv("AUTH_URL")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	parseURL := os.Getenv("PARSE_URL")

	// Authorization
	client, err := auth.Authorization(authURL, email, password)
	if err != nil {
		aws.UploadLogToS3("ProductsOne", fmt.Sprintf("Authorization error: %v", err))
	}

	// Download products list
	loadData, err := aws.DownloadJsonFromS3("ProductsOne", "bucket-products", "products.json")
	if err != nil {
		aws.UploadLogToS3("ProductsOne", fmt.Sprintf("Download products error: %v", err))
	}

	var products []shared.Product
	if err = json.Unmarshal(loadData, &products); err != nil {
		aws.UploadLogToS3("ProductsOne", fmt.Sprintf("error unmarshalling JSON data to []shared.Product: %w", err))
	}

	// Get Max Date
	maxDate := utils.FindMaxDate(products)

	// Create storage
	files.CreateFolderStorage()

	// Parse Data
	for _, market := range products {
		product := parsers.ParseDataOneProduct(client, market, parseURL)

		overviewName := fmt.Sprintf("%s.html", market.Name)

		// Write files
		err = files.WriteHTML(product, fmt.Sprintf("./storage/%s", overviewName))
		if err != nil {
			aws.UploadLogToS3("ProductsOne", fmt.Sprintf("error write file: %w", err))
		}

	}

	// Create ZIP
	err = files.CreateZipArchive("./storage", "./storage/"+maxDate+".zip")
	if err != nil {
		aws.UploadLogToS3("ProductsOne", fmt.Sprintf("Failed Create ZIP file: %w", err))
	}

	// Upload products to S3
	err = aws.UploadFileToS3("bucket-store", "./storage/"+maxDate+".zip", "daily/")
	if err != nil {
		aws.UploadLogToS3("ProductsOne", fmt.Sprintf("Failed to upload file: %w", err))
	}

	// Remove folders
	files.RemoveFolders()
}
