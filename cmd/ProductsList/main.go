package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"go-parser-lambda/pkg/auth"
	"go-parser-lambda/pkg/aws"
	"go-parser-lambda/pkg/parsers"
	"os"
)

func ProductsList(ctx context.Context) (*string, error) {
	authURL := os.Getenv("AUTH_URL")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	parseURL := os.Getenv("PARSE_URL")

	// Authorization
	client, err := auth.Authorization(authURL, email, password)
	if err != nil {
		aws.UploadLogToS3("ProductsList", fmt.Sprintf("Authorization error: %v", err))
	}

	// Get Products List
	products, maxDate := parsers.ParseProductsList(client)

	// Upload to S3
	err = aws.UploadJsonToS3("ProductsList", "bucket-products", "products.json", products)
	err = aws.UploadJsonToS3("ProductsList", "bucket-products", maxDate+".json", products)

	// Run Description Parser
	bfData := parsers.ParseDescription(client, products, parseURL)
	err = aws.UploadJsonToS3("ProductsList", "bucket-description", "product-one/"+maxDate+".json", bfData)
	return nil, nil
}

func main() {
	lambda.Start(ProductsList)
}
