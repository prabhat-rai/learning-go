package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"prabhat-rai.in/snippetbox/pkg/forms"
	"prabhat-rai.in/snippetbox/pkg/models"
	"strconv"
)

type SnippetRequest struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Expires string `json:"expires"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fmt.Printf("%#v", r.PostForm)
	form := forms.New(r.PostForm)

	if !validateCreateRequest(form, false) {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetFromApi(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	snippetRequest := SnippetRequest{}
	err = json.NewDecoder(r.Body).Decode(&snippetRequest)
	if err != nil{
		app.serverError(w, err)
		return
	}

	formValues, _ := query.Values(snippetRequest)
	// fmt.Printf("%#v", formValues)
	form := forms.New(formValues)

	if !validateCreateRequest(form, true) {
		errorsArray := form.Errors.GetAllErrors()
		errorsJson, _ := json.Marshal(errorsArray)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"errors\" : " + string(errorsJson) + "}"))
		return
	}

	id, err := app.snippets.Insert(snippetRequest.Title, snippetRequest.Content, snippetRequest.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	snippetRequest.Id = id
	//Marshal or convert user object back to json and write to response
	snippetResponse, err := json.Marshal(snippetRequest)
	if err != nil{
		app.serverError(w, err)
		return
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(snippetResponse)
}

func validateCreateRequest(form *forms.Form, fromApi bool) bool {
	titleField := "title"
	contentField := "content"
	expiresField := "expires"

	if fromApi {
		titleField = "Title"
		contentField = "Content"
		expiresField = "Expires"
	}

	form.Required(titleField, contentField, expiresField)
	form.MaxLength(titleField, 100)
	form.PermittedValues(expiresField, "365", "7", "1")

	return form.Valid()
}
