package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/wizzldev/mailer/types"
)

var listenAddr = flag.String("listenAddr", ":3000", "The address to listen on")

func main() {
	flag.Parse()
	_ = godotenv.Load()

	tmpl, err := NewTemplate("layout", "main", types.TemplateProps{
		"github_url":    os.Getenv("GITHUB_URL"),
		"discord_url":   os.Getenv("DISCORD_URL"),
		"instagram_url": os.Getenv("INSTAGRAM_URL"),
		"kofi_url":      os.Getenv("KOFI_URL"),
		"app_url":       os.Getenv("APP_URL"),
		"app_logo_url":  os.Getenv("APP_LOGO_URL"),
	})

	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	mailer := NewMailService(&Config{
		From:        os.Getenv("MAIL_NAME"),
		FromAddress: os.Getenv("MAIL_FROM"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    port,
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
	}, tmpl)

	server := NewApiServer(mailer)

	log.Printf("Server listening on %s\n", *listenAddr)
	if err := server.Start(*listenAddr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
