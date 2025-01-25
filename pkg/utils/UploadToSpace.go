package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToSpace(file io.ReadSeeker, fileKey string) error {
	SECRET := os.Getenv("SPACES_SECRET_KEY")
	key := os.Getenv("SPACES_ACCESS")

	// ContentType := "image"

	endpoint := "https://fra1.digitaloceanspaces.com"

	// Create S3 client
	
    s3Config := &aws.Config{
        Credentials: credentials.NewStaticCredentials(key, SECRET, ""), // Specifies your credentials.
        Endpoint:    aws.String(endpoint), // Find your endpoint in the control panel, under Settings. Prepend "https://".
        S3ForcePathStyle: aws.Bool(false), // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
        Region:      aws.String("fra1"), // Must be "us-east-1" when creating new Spaces. Otherwise, use the region in your endpoint, such as "nyc3".
    }

    // Step 3: The new session validates your request and directs it to your Space's specified endpoint using the AWS SDK.
    newSession, err := session.NewSession(s3Config)
    if err != nil {
        return err
    }

    s3Client := s3.New(newSession)

	// Detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	file.Seek(0, io.SeekStart) // Reset file pointer
	contentType := http.DetectContentType(buffer)

	object := s3.PutObjectInput{
		Bucket:      aws.String("ddl"),
		Key:         aws.String(fileKey),
		Body:        file,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	}


	_, err = s3Client.PutObject(&object) 
	if err != nil { 
		return err
	}

	return nil

}