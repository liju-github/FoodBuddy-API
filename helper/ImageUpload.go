package helper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func ImageUpload(c *gin.Context) {

	//CLOUDINARY_URL=cloudinary://<api_key>:<api_secret>@<cloud_name>

	url := fmt.Sprintf("cloudinary://%v:%v@%v", GetEnvVariables().CloudinaryAccessKey, GetEnvVariables().CloudinarySecretKey, GetEnvVariables().CloudinaryCloudName)
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		log.Printf("Failed to initialize Cloudinary: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
		return
	}

	ctx := context.Background()

	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Failed to get multipart form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
		return
	}

	files := form.File["files"]
	var uploadedURLs []string

	for _, fileHeader := range files {
		if !isValidFormat(fileHeader.Filename) {
			log.Printf("Invalid file format: %s", fileHeader.Filename)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		// read the file into a buffer
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Printf("Failed to read file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		// define the upload parameters with the desired transformations
		uploadParams := uploader.UploadParams{
			Transformation: "f_auto/q_auto/c_crop,w_300,h_300,c_fill", 
			Folder:         "foodbuddy",   
		}

		// upload the file with transformation to Cloudinary
		uploadResult, err := cld.Upload.Upload(ctx, buf, uploadParams)
		if err != nil {
			log.Printf("Failed to upload file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		fmt.Println(uploadResult)
		uploadedURLs = append(uploadedURLs, uploadResult.SecureURL)
	}

	c.JSON(http.StatusOK, gin.H{"uploaded_urls": uploadedURLs})

}

//to check the file type 
func isValidFormat(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
