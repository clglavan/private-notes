package privateNotes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// func init() {
// 	// Register an HTTP function with the Functions Framework
// 	functions.HTTP("privateNotes", PrivateNotes)
// }

type SecretNote struct {
	Key               string
	SecureNote        string
	RecaptchaResponse string
}

type IndexPageData struct {
	PostUrl string
}
type ConfirmPageData struct {
	PostUrl string
	Key     string
}

type SuccessPageData struct {
	SecretUrl string
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
}

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

func PrivateNotes(w http.ResponseWriter, r *http.Request) {

	PUBLIC_URL := os.Getenv("PUBLIC_URL")

	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	DEFAULT_EXPIRATION := os.Getenv("DEFAULT_EXPIRATION")
	DEFAULT_EXPIRATION_INT, err := strconv.Atoi(DEFAULT_EXPIRATION)
	RECAPTCHA_SECRET := os.Getenv("RECAPTCHA_SECRET")

	if err != nil {
		fmt.Println("Default expiration is not an integer")
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: REDIS_PASSWORD,
		DB:       0,
	})

	secretURL := "?key="

	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if key != "" {
			data := ConfirmPageData{
				PostUrl: PUBLIC_URL,
				Key:     key,
			}
			tmpl := template.Must(template.ParseFiles("views/layout.html", "views/confirm.html"))
			tmpl.ParseGlob("views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		} else {
			data := IndexPageData{
				PostUrl: PUBLIC_URL,
			}
			tmpl := template.Must(template.ParseFiles("views/layout.html", "views/index.html"))
			tmpl.ParseGlob("views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		}
	case http.MethodPost:
		function := r.FormValue("function")
		switch function {
		case "create":
			// ##################### Get the form data
			r.ParseForm()
			var t SecretNote
			t.Key = r.FormValue("key")
			t.SecureNote = r.FormValue("secureNote")
			t.RecaptchaResponse = r.FormValue("g-recaptcha-response")

			// Check and verify the recaptcha response token.
			if err := CheckRecaptcha(RECAPTCHA_SECRET, t.RecaptchaResponse); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// ##################### Prepare the url
			data := SuccessPageData{
				SecretUrl: string(secretURL + t.Key),
			}
			// ##################### Save the cipherText to redis

			ctx := context.Background()

			err := rdb.Set(ctx, t.Key, t.SecureNote, time.Second*time.Duration(DEFAULT_EXPIRATION_INT)).Err()
			if err != nil {
				panic(err)
			}

			// ##################### Render the reponse template
			tmpl := template.Must(template.ParseFiles("views/layout.html", "views/success.html"))
			tmpl.ParseGlob("views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		case "retrieve":
			key := r.FormValue("key")
			if key != "" {
				ctx := context.Background()

				val, err := rdb.Get(ctx, key).Result()

				if err != nil {
					tmpl := template.Must(template.ParseFiles("views/layout.html", "views/error.html"))
					tmpl.ParseGlob("views/assets/*")
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					tmpl.Execute(w, "")
					return
				}

				data := SecretNote{
					Key:        key,
					SecureNote: string(val),
				}
				tmpl := template.Must(template.ParseFiles("views/layout.html", "views/result.html"))
				tmpl.ParseGlob("views/assets/*")
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				tmpl.Execute(w, data)

				rdb.Del(ctx, key)

			}

		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CheckRecaptcha(secret, response string) error {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	// Check recaptcha verification success.
	if !body.Success {
		return errors.New("unsuccessful recaptcha verify request")
	}

	return nil
}
