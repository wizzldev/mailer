package client

// Message represents a message to be sent.
type Message struct {
	ToAddress string    `json:"address"`
	Subject   string    `json:"subject"`
	Content   *Content  `json:"content,omitempty"`
	Template  *Template `json:"template,omitempty"`
}

// Template represents a template to be used for a message.
type Template struct {
	ID    string        `json:"id"`
	Props TemplateProps `json:"props,omitempty"`
}

// Content represents the content of a message.
type Content struct {
	Body string `json:"body"`
	HTML bool   `json:"html,omitempty"`
}

// Templates

// TemplateProp represents a map of string keys and any type of values.
type TemplateProps map[string]any

// TemplateID represents a unique identifier for a template.
type TemplateID string

const (
	TemplateRegister       TemplateID = "register"
	TemplateForgotPassword TemplateID = "forgot_password"
)

func (id TemplateID) String() string {
	return string(id)
}
