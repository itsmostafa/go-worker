package aws

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESEmail struct {
	From        string
	To          string
	Subject     string
	CC          string `default:""`
	BCC         string `default:""`
	HtmlBody    string
	TextBody    string
	PdfFileName string
	PdfFile     []byte
}

func buildMessage(source, destination, cc, bcc, subject, html, text, pdfFileName string,
	pdfFile []byte) (*ses.SendRawEmailInput, error) {

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// Email Header
	h := make(textproto.MIMEHeader)
	h.Set("From", source)
	h.Set("To", destination)
	h.Set("CC", cc)
	h.Set("Return-Path", source)
	h.Set("Subject", subject)
	h.Set("Content-Language", "en-US")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// Text body
	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "8bit")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(text))
	if err != nil {
		return nil, err
	}

	// HTML body
	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	h.Set("Content-Type", "text/html; charset=utf-8")
	part, err = writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(html))
	if err != nil {
		return nil, err
	}

	// File Attachment
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment; filename="+pdfFileName)
	h.Set("Content-Type", "application/pdf; x-unix-mode=0644; name=\""+pdfFileName+"\"")
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	part, err = writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(pdfFile)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Strip boundary line before header (doesn't work with it present)
	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		return nil, fmt.Errorf("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]

	raw := ses.RawMessage{
		Data: []byte(s),
	}
	input := &ses.SendRawEmailInput{
		Destinations: []*string{aws.String(destination)},
		Source:       aws.String(source),
		RawMessage:   &raw,
	}

	return input, nil
}

func (e *SESEmail) Send() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		return err
	}

	input, err := buildMessage(e.From, e.To, e.CC, e.BCC, e.Subject, e.HtmlBody, e.TextBody, e.PdfFileName, e.PdfFile)
	if err != nil {
		return err
	}

	svc := ses.New(sess)

	result, err := svc.SendRawEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	fmt.Println(result)
	return nil
}
