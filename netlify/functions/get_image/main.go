package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const baseURL = "https://dchqsxvlazultvsphmau.supabase.co/storage/v1/object/public/"

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	imagePath := strings.TrimPrefix(request.Path, "/.netlify/functions/get_image/")
	imageURL := baseURL + imagePath

	// Fetch the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &events.APIGatewayProxyResponse{
			StatusCode: resp.StatusCode,
			Body:       fmt.Sprintf("Failed to get the image: %s", resp.Status),
		}, nil
	}

	// Read the image
	imgBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	// Set the correct Content-Type based on your image format
	contentType := http.DetectContentType(imgBytes)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		Body:            base64.StdEncoding.EncodeToString(imgBytes),
		IsBase64Encoded: true,
	}, nil
}

func main() {
	lambda.Start(handler)
}