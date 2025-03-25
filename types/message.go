package types

// Message represents a message to be sent.
type Message struct {
	ToAddress string    `json:"address"`
	Subject   string    `json:"subject"`
	Content   *Content  `json:"content,omitempty"`
	Template  *Template `json:"template,omitempty"`
}

type Template struct {
	ID    string        `json:"id"`
	Props TemplateProps `json:"props,omitempty"`
}

type Content struct {
	Body string `json:"body"`
	HTML bool   `json:"html,omitempty"`
}
