package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	endpoint   string
	bucketName string
	region     string
)

func init() {
	endpoint = "http://localhost:4572/"
	bucketName = "test-bucket/"
	region = "us-east-1"
}

//RequestPDF struct
type RequestPDF struct {
	body []byte
}

// Data struct, collection to append data to html
type Data struct {
	Product       []Product
	Image         string
	InvoiceNo     int
	Date          string
	DueDate       string
	DeliveryDate  string
	PaymentMethod string
}

// Product struct, collection to list order product
type Product struct {
	Item     string
	Price    float64
	Qty      float64
	Subtotal float64
}

func main() {
	// fill variable in dynamic page template
	var (
		products []Product
	)
	dt := time.Now()
	data := Data{
		Image:         "https://camo.githubusercontent.com/27a12a480bd0a8c936bc6cf3edf2ab0fcbd4dd38/68747470733a2f2f7261772e6769746875622e636f6d2f676f6c616e672d73616d706c65732f676f706865722d766563746f722f6d61737465722f676f706865722d736964655f706174682e706e67",
		InvoiceNo:     1,
		Date:          dt.Format("01-02-2006"),
		DueDate:       dt.Format("01-02-2006"),
		DeliveryDate:  dt.Format("01-02-2006"),
		PaymentMethod: "Paypal",
	}
	products = append(products, Product{Item: "item", Price: 10000, Qty: 1, Subtotal: 10000})
	products = append(products, Product{Item: "item2", Price: 20000, Qty: 2, Subtotal: 40000})
	products = append(products, Product{Item: "item3", Price: 30000, Qty: 3, Subtotal: 90000})

	data.Product = products
	r := &RequestPDF{}
	err := r.ParseTemplate("invoice.html", data)

	// create pdf from html
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(r.body)))
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile("./test.pdf")
	// if err != nil {
	// log.Fatal(err)
	// }

	// process sending to aws s3
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(
			"foo", // id
			"foo", // secret
			""),   // token can be left blank for now
	})

	if err != nil {
		log.Fatal(err)
	}

	// Upload
	fileName, err := AddFileToS3(s, pdfg.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done with url file : %s%s%s", endpoint, bucketName, fileName)
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, buffer []byte) (string, error) {
	fileName := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(len(buffer))),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return fileName, err
}

// ParseTemplate HTML from file
func (r *RequestPDF) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.Bytes()
	return nil
}
