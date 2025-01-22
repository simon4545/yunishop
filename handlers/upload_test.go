package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestUploadImage(t *testing.T) {
	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add a file to the request
	file, err := os.Open("huwai.jpg")
	if err != nil {
		t.Fatalf("Failed to open test image: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "huwai.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file content: %v", err)
	}

	// Add the product_id field to the request
	err = writer.WriteField("product_id", "12345")
	if err != nil {
		t.Fatalf("Failed to write product_id field: %v", err)
	}

	// Close the writer to finalize the multipart message
	writer.Close()

	// Create a new request with the multipart body
	req, err := http.NewRequest("POST", "/images", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rr)
	// Call the UploadImage handler
	UploadImage(c)

	// Check the status codec
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message":"Image uploaded successfully"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
