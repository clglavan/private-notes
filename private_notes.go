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
	CUSTOM_LOGO       string
	expiration        time.Duration
	Lang              LangData
}

type IndexPageData struct {
	PostUrl                string
	DEFAULT_EXPIRATION_INT int
	MAXIMUM_EXPIRATION_INT int
	ErrorBag               []string
	CUSTOM_LOGO            string
	Lang                   LangData
}
type ConfirmPageData struct {
	PostUrl     string
	Key         string
	CUSTOM_LOGO string
	Lang        LangData
}

type ErrotPageData struct {
	CUSTOM_LOGO string
	Lang        LangData
}
type SuccessPageData struct {
	SecretUrl   string
	CUSTOM_LOGO string
	Lang        LangData
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
}

type LangData struct {
	INDEX_TITLE                       string
	INDEX_SUBTITLE                    string
	INDEX_NOTE_PLACEHOLDER            string
	INDEX_PASSWORD                    string
	INDEX_PASSWORD_PLACEHOLDER        string
	INDEX_EXPIRATION                  string
	INDEX_SEND_BUTTON                 string
	SUCCESS_TITLE                     string
	SUCCESS_SUBTITLE                  string
	SUCCESS_TOOLTIP                   string
	CONFIRM_SUBTITLE                  string
	CONFIRM_SHOW_BUTTON               string
	RESULT_TITLE                      string
	RESULT_SUBTITLE                   string
	RESULT_PASSWORD                   string
	RESULT_PASSWORD_PLACEHOLDER       string
	RESULT_TOOLTIP                    string
	ERROR_TITLE                       string
	ERROR_SUBTITLE                    string
	LANG_ERRORBAG_EMPTY               string
	LANG_ERRORBAG_PASSWORD_REQUIRED   string
	LANG_ERRORBAG_EXPIRATION_REQUIRED string
	LANG_ERRORBAG_NOTE_TOO_LONG       string
	LANG_ERRORBAG_EXPIRATION_TOO_LONG string
}

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

func PrivateNotes(w http.ResponseWriter, r *http.Request) {

	lang := LangData{
		INDEX_TITLE:                       os.Getenv("LANG_INDEX_TITLE"),
		INDEX_SUBTITLE:                    os.Getenv("LANG_INDEX_SUBTITLE"),
		INDEX_NOTE_PLACEHOLDER:            os.Getenv("LANG_INDEX_NOTE_PLACEHOLDER"),
		INDEX_PASSWORD:                    os.Getenv("LANG_INDEX_PASSWORD"),
		INDEX_PASSWORD_PLACEHOLDER:        os.Getenv("LANG_INDEX_PASSWORD_PLACEHOLDER"),
		INDEX_EXPIRATION:                  os.Getenv("LANG_INDEX_EXPIRATION"),
		INDEX_SEND_BUTTON:                 os.Getenv("LANG_INDEX_SEND_BUTTON"),
		SUCCESS_TITLE:                     os.Getenv("LANG_SUCCESS_TITLE"),
		SUCCESS_SUBTITLE:                  os.Getenv("LANG_SUCCESS_SUBTITLE"),
		SUCCESS_TOOLTIP:                   os.Getenv("LANG_SUCCESS_TOOLTIP"),
		CONFIRM_SUBTITLE:                  os.Getenv("LANG_CONFIRM_SUBTITLE"),
		CONFIRM_SHOW_BUTTON:               os.Getenv("LANG_CONFIRM_SHOW_BUTTON"),
		RESULT_TITLE:                      os.Getenv("LANG_RESULT_TITLE"),
		RESULT_SUBTITLE:                   os.Getenv("LANG_RESULT_SUBTITLE"),
		RESULT_PASSWORD:                   os.Getenv("LANG_RESULT_PASSWORD"),
		RESULT_PASSWORD_PLACEHOLDER:       os.Getenv("LANG_RESULT_PASSWORD_PLACEHOLDER"),
		RESULT_TOOLTIP:                    os.Getenv("LANG_RESULT_TOOLTIP"),
		ERROR_TITLE:                       os.Getenv("LANG_ERROR_TITLE"),
		ERROR_SUBTITLE:                    os.Getenv("LANG_ERROR_SUBTITLE"),
		LANG_ERRORBAG_EMPTY:               os.Getenv("LANG_ERRORBAG_EMPTY"),
		LANG_ERRORBAG_PASSWORD_REQUIRED:   os.Getenv("LANG_ERRORBAG_PASSWORD_REQUIRED"),
		LANG_ERRORBAG_EXPIRATION_REQUIRED: os.Getenv("LANG_ERRORBAG_EXPIRATION_REQUIRED"),
		LANG_ERRORBAG_NOTE_TOO_LONG:       os.Getenv("LANG_ERRORBAG_NOTE_TOO_LONG"),
		LANG_ERRORBAG_EXPIRATION_TOO_LONG: os.Getenv("LANG_ERRORBAG_EXPIRATION_TOO_LONG"),
	}

	PUBLIC_URL := os.Getenv("PUBLIC_URL")
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	DEFAULT_EXPIRATION := os.Getenv("DEFAULT_EXPIRATION")
	DEFAULT_EXPIRATION_INT, err := strconv.Atoi(DEFAULT_EXPIRATION)
	RECAPTCHA_SECRET := os.Getenv("RECAPTCHA_SECRET")
	CUSTOM_LOGO := os.Getenv("CUSTOM_LOGO")

	if err != nil {
		fmt.Println("Default expiration is not an integer")
		return
	}

	MAXIMUM_EXPIRATION := os.Getenv("MAXIMUM_EXPIRATION")
	MAXIMUM_EXPIRATION_INT, err := strconv.Atoi(MAXIMUM_EXPIRATION)
	if err != nil {
		fmt.Println("Maximum expiration is not an integer")
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
				PostUrl:     PUBLIC_URL,
				CUSTOM_LOGO: CUSTOM_LOGO,
				Key:         key,
				Lang:        lang,
			}
			tmpl := template.Must(template.ParseFiles("views/layout.html", "views/confirm.html"))
			tmpl.ParseGlob("views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		} else {
			data := IndexPageData{
				PostUrl:                PUBLIC_URL,
				DEFAULT_EXPIRATION_INT: DEFAULT_EXPIRATION_INT / 60,
				MAXIMUM_EXPIRATION_INT: MAXIMUM_EXPIRATION_INT / 60,
				CUSTOM_LOGO:            CUSTOM_LOGO,
				ErrorBag:               nil,
				Lang:                   lang,
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

			if r.FormValue("expirationTime") != "" {

				intExpiration, err := strconv.Atoi(r.FormValue("expirationTime"))
				if err != nil {
					fmt.Println("Default expiration is not an integer")
					return
				}

				if intExpiration*60 > MAXIMUM_EXPIRATION_INT {
					data := IndexPageData{
						PostUrl:                PUBLIC_URL,
						DEFAULT_EXPIRATION_INT: DEFAULT_EXPIRATION_INT / 60,
						CUSTOM_LOGO:            CUSTOM_LOGO,
						Lang:                   lang,
						ErrorBag:               []string{"Failed! The expiration amount exceeds the maximum of " + strconv.Itoa(MAXIMUM_EXPIRATION_INT/60) + " minutes"},
					}
					tmpl := template.Must(template.ParseFiles("views/layout.html", "views/index.html"))
					tmpl.ParseGlob("views/assets/*")
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					tmpl.Execute(w, data)
					return
				} else {
					t.expiration = time.Second * time.Duration(intExpiration*60)
				}
			} else {
				t.expiration = time.Second * time.Duration(DEFAULT_EXPIRATION_INT)
			}

			// Check and verify the recaptcha response token.
			if err := CheckRecaptcha(RECAPTCHA_SECRET, t.RecaptchaResponse); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// ##################### Prepare the url
			data := SuccessPageData{
				SecretUrl:   string(secretURL + t.Key),
				CUSTOM_LOGO: CUSTOM_LOGO,
				Lang:        lang,
			}
			// ##################### Save the cipherText to redis

			ctx := context.Background()

			err := rdb.Set(ctx, t.Key, t.SecureNote, t.expiration).Err()
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
					data := ErrotPageData{
						CUSTOM_LOGO: CUSTOM_LOGO,
						Lang:        lang,
					}
					tmpl := template.Must(template.ParseFiles("views/layout.html", "views/error.html"))
					tmpl.ParseGlob("views/assets/*")
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					tmpl.Execute(w, data)
					return
				}

				data := SecretNote{
					Key:         key,
					CUSTOM_LOGO: CUSTOM_LOGO,
					SecureNote:  string(val),
					Lang:        lang,
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
