package privateNotes

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPrivateNotes(t *testing.T) {

	t.Log("Preparing the environment")
	t.Setenv("RECAPTCHA_SECRET", "6LfrnPQkAAAAAIPBpB3B_SfK5oC-2BIHzkI29ozF")
	t.Setenv("RECAPTCHA_KEY", "6LfrnPQkAAAAAKSC7plAfp8kutlL9cKpvSWsqgiL")
	t.Setenv("REDIS_HOST", "redis")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("DEFAULT_EXPIRATION", "3600")
	t.Setenv("MAXIMUM_EXPIRATION", "3600")
	t.Setenv("CUSTOM_LOGO", "https://raw.githubusercontent.com/clglavan/private-notes/add_logo/logo.png")
	t.Setenv("NOTE_MAX_LENGTH_SERVER", "20000")
	t.Setenv("NOTE_MAX_LENGTH_CLIENT", "10000")
	t.Setenv("LANG_INDEX_TITLE", "Private notes")
	t.Setenv("LANG_INDEX_SUBTITLE", "Send self-destructing private notes securely")
	t.Setenv("LANG_INDEX_NOTE_PLACEHOLDER", "Notes will self-destruct after they are read...")
	t.Setenv("LANG_INDEX_PASSWORD", "Password protect your note")
	t.Setenv("LANG_INDEX_PASSWORD_PLACEHOLDER", "enter your password")
	t.Setenv("LANG_INDEX_EXPIRATION", "Custom expiration time,1-60,default 60 min")
	t.Setenv("LANG_INDEX_SEND_BUTTON", "Encrypt & Send")
	t.Setenv("LANG_SUCCESS_TITLE", "Thank you for using private notes")
	t.Setenv("LANG_SUCCESS_SUBTITLE", "This note will self-destruct after it will be read. Click on it to copy to clipboard and send this link to the other party.")
	t.Setenv("LANG_SUCCESS_TOOLTIP", "Click on it to copy to clipboard")
	t.Setenv("LANG_CONFIRM_SUBTITLE", "Do you want to decrypt this message now? It's contents will be lost forever")
	t.Setenv("LANG_CONFIRM_SHOW_BUTTON", "Show & Destroy")
	t.Setenv("LANG_RESULT_TITLE", "Thank you for using private notes")
	t.Setenv("LANG_RESULT_SUBTITLE", "This note has been destroyed, below is the only copy.")
	t.Setenv("LANG_RESULT_PASSWORD", "This message is password protected")
	t.Setenv("LANG_RESULT_PASSWORD_PLACEHOLDER", "Enter your password")
	t.Setenv("LANG_RESULT_TOOLTIP", "Click on it to copy to clipboard")
	t.Setenv("LANG_ERROR_TITLE", "Private notes")
	t.Setenv("LANG_ERROR_SUBTITLE", "Note does not exist")
	t.Setenv("LANG_ERRORBAG_EMPTY", "You can't send an empty note")
	t.Setenv("LANG_ERRORBAG_PASSWORD_REQUIRED", "Secret password checked but not provided")
	t.Setenv("LANG_ERRORBAG_EXPIRATION_REQUIRED", "Expiration time checked but not provided")
	t.Setenv("LANG_ERRORBAG_NOTE_TOO_LONG", "Secret note is too long")
	t.Setenv("LANG_ERRORBAG_EXPIRATION_TOO_LONG", "Expiration time is too long")

	// Validate the GET / response 1)
	t.Log("1) Validate the GET / response")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	PrivateNotes(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("1.expected error to be nil got %v", err)
	}
	// 1.a Check for textarea
	t.Log("1.a Checkin for the textarea for input")
	if !strings.Contains(string(data), "<textarea id=\"secureNote\"") {
		t.Errorf("1.Expected a textarea but wasn't found, got %v", string(data))
	}
	// 1.b Check for button
	t.Log("1.b Checking for the submit button")
	if !strings.Contains(string(data), "<input type=\"submit\"") {
		t.Errorf("Expected a submit button but wasn't found, got %v", string(data))
	}

	/* ============================================================================= */
	// Validate the POST / response 2)
	formData := url.Values{}
	formData.Set("function", "create")
	formData.Add("secureNote", "468cdaed482145453bbbceb74629633951fb10303c15c55b66a367b54716aaa1xGlMEy3u/f2UzQpKtYANuw==45b517b96c60b8a13628eb51291c32856dfba3b5263e466e156b3cd3b5b70111")
	formData.Add("secretPassword", "")
	formData.Add("expirationTime", "")
	formData.Add("key", "55a3932290cb72fbc28f5682b4da1e7e2c0c18223a28746ab6953a87b5013f8d")
	formData.Add("key", "55a3932290cb72fbc28f5682b4da1e7e2c0c18223a28746ab6953a87b5013f8d")
	formData.Add("g-recaptcha-response", "")

	t.Log("2) Validate the POST / response")
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	w = httptest.NewRecorder()
	PrivateNotes(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("2.expected error to be nil got %v", err)
	}
	// 2.a Check for secretUrl
	t.Log("2.a Checking for secreUrl link")
	if !strings.Contains(string(data), "id=\"secretURL\"") {
		t.Errorf("2.a Expected a secreUrl link but wasn't found, got %v", string(data))
	}
	// // 2.b Check for button
	// t.Log("2.b Checking for the submit button")
	// if !strings.Contains(string(data), "<input type=\"submit\"") {
	// 	t.Errorf("2.b Expected a submit button but wasn't found, got %v", string(data))
	// }

}
