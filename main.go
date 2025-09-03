package main

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/character"
	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func mustGenerateCSRFKey() (key []byte) {
	key = make([]byte, 32)
	n, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	if n != 32 {
		panic("unable to read 32 bytes for CSRF key")
	}
	return
}

var textToDraw chan DisplayUpdate = make(chan DisplayUpdate, 1000000)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	_, err := host.Init()
	if err != nil {
		log.Error("error initializing periph", slog.Any("error", err))
		return
	}
	bus, err := i2creg.Open("")
	if err != nil {
		log.Error("error opening i2c", slog.Any("error", err))
		return
	}
	dev := &i2c.Dev{
		Bus:  bus,
		Addr: 0x27,
	}
	d := character.NewDisplay(dev, false)
	d.SetBacklight(false)
	d.Clear()

	r := http.NewServeMux()
	r.Handle("/", http.HandlerFunc(handler))
	csrfMiddleware := csrf.Protect(mustGenerateCSRFKey(), csrf.TrustedOrigins([]string{"lcd.edwardh.dev"}))
	withCSRFProtection := csrfMiddleware(r)

	lastMessageTime := time.Now()

	go func() {
		textToDraw <- DisplayUpdate{Line1: "lcd.edwardh.dev", Line2: ""}
		for txt := range textToDraw {
			d.SetBacklight(true)
			d.Clear()
			d.Goto(0, 0)
			d.Print(txt.Line1)
			d.Goto(1, 0)
			d.Print(txt.Line2)
			time.Sleep(250 * time.Millisecond)
			lastMessageTime = time.Now()
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if time.Since(lastMessageTime) > 30*time.Second {
				d.SetBacklight(false)
				d.Goto(0, 0) // backlight does not change until another command is sent
			}
			if time.Since(lastMessageTime) > 10*time.Minute {
				d.Clear()
			}
		}
	}()

	if err := http.ListenAndServe(":8019", withCSRFProtection); err != nil {
		log.Error("error starting web server", slog.Any("error", err))
		return
	}
}

type DisplayUpdate struct {
	Line1 string
	Line2 string
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templ.Handler(page(form()), templ.WithStreaming()).ServeHTTP(w, r)
	case http.MethodPost:
		UpdateHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	line1 := sanitizeInput(r.Form.Get("line1"))
	line2 := sanitizeInput(r.Form.Get("line2"))
	select {
	case textToDraw <- DisplayUpdate{Line1: line1, Line2: line2}:
		templ.Handler(page(success()), templ.WithStreaming()).ServeHTTP(w, r)
	default:
		http.Error(w, "too many requests", http.StatusTooManyRequests)
	}
}

// limit to ASCII and truncate to 16 characters
func sanitizeInput(input string) string {
	var sanitized string
	for _, r := range input {
		if r < 32 || r > 126 {
			continue
		}
		sanitized += string(r)
		if len(sanitized) >= 16 {
			break
		}
	}
	return sanitized
}
