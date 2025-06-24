package main

import (
	"fmt"
	//"html/template"
	"net/http"
	"strconv"
	"errors"

	"GoBin/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets,err := app.snippets.Latest()
	if err != nil {
		app.serverError(w,err)
	}

	for _,snippet := range snippets {
		fmt.Fprintf(w,"%+v\n",snippet)
	}

	/* files := []string{"./ui/html/base.tmpl", "./ui/html/pages/home.tmpl", "./ui/html/partials/nav.tmpl"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} */
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet,err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err,models.ErrNoRecords){
			app.notFound(w)
		} else {
			app.serverError(w,err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id,err := app.snippets.Insert(title,content,expires)
	if err != nil {
		app.serverError(w,err)
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view?id=%v",id),http.StatusSeeOther)
}
