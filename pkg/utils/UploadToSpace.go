package utils

import (
	"io"
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

	// FOR SHAREPICS ONLY
    // Step 4: Define the parameters of the object you want to upload.
    object := s3.PutObjectInput{
        Bucket: aws.String("ddl"), // The path to the directory you want to upload the object to, starting with your Space name.
        Key:    aws.String(fileKey), // Object key, referenced whenever you want to access this file later.
        Body:   file, // The object's contents.
        ACL:    aws.String("public-read"), // Defines Access-control List (ACL) permissions, such as private or public.
        // Metadata: map[string]*string{ // Required. Defines metadata tags.
        //                         "x-amz-meta-my-key": aws.String("DEBUG"),
        //                 },
    }

	_, err = s3Client.PutObject(&object) 
	if err != nil { 
		return err
	}

	return nil

}