package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wizzldev/mailer/types"
)

// ApiServer represents an API server.
type ApiServer struct {
	svc Service
}

// NewApiServer creates a new API server.
func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

// Start starts the API server.
func (s *ApiServer) Start(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s)
}

// ServeHTTP handles incoming HTTP requests.
func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check if the request path is /send
	if r.URL.Path != "/send" {
		writeError(w, "This resource is not found.", http.StatusNotFound)
		return
	}

	// check if the request method is POST
	if r.Method != http.MethodPost {
		writeError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// read the request body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data types.Message
	if err := json.Unmarshal(b, &data); err != nil {
		writeError(w, fmt.Sprintf("Failed to unmarshal request body: %v", err), http.StatusBadRequest)
		return
	}

	// initialize the message init struct
	init := MessageInit{
		ToAddress: data.ToAddress,
		Subject:   data.Subject,
	}

	// if the message is a template
	if data.Template != nil {
		if err := s.svc.SendTemplate(init, data.Template.ID, data.Template.Props); err != nil {
			writeError(w, fmt.Sprintf("Failed to send template email: %v", err), http.StatusInternalServerError)
			return
		}

		// if the message is a template, we don't need to send it again
		writeSuccess(w)
		return
	}

	// write an error if both content and template are nil
	if data.Content == nil {
		writeError(w, fmt.Sprintf("No content or template provided"), http.StatusBadRequest)
		return
	}

	// create a sender function based on the content type
	var sender = s.svc.Send
	if data.Content.HTML {
		sender = s.svc.SendHTML
	}

	// send the email
	if err := sender(init, data.Content.Body); err != nil {
		writeError(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	// write a success response to the client
	writeSuccess(w)
}

// writeSuccess writes a success response to the client.
func writeSuccess(w http.ResponseWriter) {
	log.Println("Email Sent")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Email sent successfully",
	})
}

// writeError writes an error response to the client.
func writeError(w http.ResponseWriter, message string, statusCode int) {
	log.Println(message, statusCode)
	http.Error(w, message, statusCode)
}
