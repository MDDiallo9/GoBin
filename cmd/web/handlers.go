package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
	//"html/template"
	"errors"
	"net/http"
	"strconv"

	"GoBin/internal/models"

	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w,http.StatusOK,"home.tmpl",data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	if r.URL.Query().Get("dl") != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=snippet-%v.txt", snippet.Title))
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.Write([]byte(snippet.Content))
        return
	}
	
	

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w,http.StatusOK,"view.tmpl",data)
	
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w,http.StatusOK,"create.tmpl",data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w,http.StatusBadRequest)
	}

	expires,err := strconv.Atoi(r.PostForm.Get("expires"))
	form := snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content:r.PostForm.Get("content"),
		Expires: expires,
		FieldErrors: make(map[string]string),
	}

	if err != nil {
		app.clientError(w,http.StatusBadRequest)
	}

	// Checking inputs

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 { // better than len(title) because len only counts bytes , this one counts runes
		form.FieldErrors["title"] = "This field cannot be longer than a 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "This field must be either 1,7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w,http.StatusOK,"create.tmpl",data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)
}

func (app *application) snippetDelete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	del,err := app.snippets.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Printf("Deleted snippet #%v\n",del)
	http.Redirect(w,r,"/",http.StatusSeeOther)
}
