package email

type Email struct {
	Sender    string
	Recipient string
	Subject   string
	HtmlBody  string
	TextBody  string
	CharSet   string
}

func (e *Email) BuildMessage() error {
	return nil
}

func (e *Email) Send() error {
	return nil
}
