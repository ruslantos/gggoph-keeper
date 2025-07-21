package post_credentials

type SavePasswordRequest struct {
	Resource string `json:"resource"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SaveCardRequest struct {
	Bank       string `json:"bank"`
	PAN        string `json:"pan"`
	ValidThru  string `json:"validThru"`
	Cardholder string `json:"cardholder"`
}

type SavePlainTextRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CredsRequest struct {
	Passwords *SavePasswordRequest  `json:"passwords,omitempty"`
	Cards     *SaveCardRequest      `json:"cards,omitempty"`
	PlainText *SavePlainTextRequest `json:"plainText,omitempty"`
}
