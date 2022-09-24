package privateNotes

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	// Register an HTTP function with the Functions Framework
	functions.HTTP("privateNotes", PrivateNotes)
}

type SecretNote struct {
	Key        string
	SecureNote string
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

func PrivateNotes(w http.ResponseWriter, r *http.Request) {

	// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./key.json") // for local development
	ENVIRONMENT := os.Getenv("ENVIRONMENT")
	function_path := "./"
	if ENVIRONMENT == "cloudfunction" {
		function_path = "./serverless_function_source_code/"
	}

	GCP_PROJECT := os.Getenv("GCP_PROJECT")
	GCP_REGION := os.Getenv("GCP_REGION")
	GCP_CF_NAME := os.Getenv("GCP_CF_NAME")
	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if key != "" {
			data := ConfirmPageData{
				PostUrl: GCP_CF_NAME,
				Key:     key,
			}
			tmpl := template.Must(template.ParseFiles(function_path+"views/layout.html", function_path+"views/confirm.html"))
			tmpl.ParseGlob(function_path + "views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		} else {
			data := IndexPageData{
				PostUrl: GCP_CF_NAME,
			}
			tmpl := template.Must(template.ParseFiles(function_path+"views/layout.html", function_path+"views/index.html"))
			tmpl.ParseGlob(function_path + "views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		}
	case http.MethodPost:
		function := r.FormValue("function")
		switch function {
		case "create":
			log.Printf("post")
			// ##################### Get the form data
			r.ParseForm()
			var t SecretNote
			t.Key = r.FormValue("key")
			t.SecureNote = r.FormValue("secureNote")
			// log.Println(t.Test)
			log.Println(t.Key)
			log.Println(t.SecureNote)
			// ##################### Prepare the url
			data := SuccessPageData{
				SecretUrl: string(GCP_REGION + "-" + GCP_PROJECT + ".cloudfunctions.net/" + GCP_CF_NAME + "?key=" + t.Key),
			}
			// ##################### Save the cipherText to bucket
			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			wc := client.Bucket("private-notes").Object(t.Key).NewWriter(ctx)
			wc.ContentType = "text/plain"
			// wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
			if _, err := wc.Write([]byte(t.SecureNote)); err != nil {
				// TODO: handle error.
				// Note that Write may return nil in some error situations,
				// so always check the error from Close.
				fmt.Println("Error: ", err)
			}
			if err := wc.Close(); err != nil {
				fmt.Println("Error: ", err)
			}

			// ##################### Render the reponse template
			tmpl := template.Must(template.ParseFiles(function_path+"views/layout.html", function_path+"views/success.html"))
			tmpl.ParseGlob(function_path + "views/assets/*")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, data)
			return
		case "retrieve":
			key := r.FormValue("key")
			log.Printf("ok")
			if key != "" {
				log.Printf("all good")
				ctx := context.Background()
				client, err := storage.NewClient(ctx)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				rc, err := client.Bucket("private-notes").Object(key).NewReader(ctx)
				if err != nil {
					fmt.Println("Error: ", err)
					// http.Error(w, "Note does not exist", http.StatusNotFound)
					tmpl := template.Must(template.ParseFiles(function_path+"views/layout.html", function_path+"views/error.html"))
					tmpl.ParseGlob(function_path + "views/assets/*")
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					tmpl.Execute(w, "")
					return
				}
				// defer
				slurp, err := ioutil.ReadAll(rc)
				rc.Close()
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				fmt.Println(string(slurp))

				if err := client.Bucket("private-notes").Object(key).Delete(ctx); err != nil {
					fmt.Println("Error: ", err)
				}

				data := SecretNote{
					Key:        key,
					SecureNote: string(slurp),
				}
				tmpl := template.Must(template.ParseFiles(function_path+"views/layout.html", function_path+"views/result.html"))
				tmpl.ParseGlob(function_path + "views/assets/*")
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				tmpl.Execute(w, data)

			}

		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
