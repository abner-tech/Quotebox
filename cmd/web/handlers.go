package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"amencia.net/quotebox/pkg/models"
)

// A struct to hold a quote
type Quotation struct {
	Quotations_id  int
	Insertion_date time.Time
	Author_name    string
	Category       string
	Quote          string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	q, err := app.quotes.Read()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Quotes: q,
	}

	//Display Quotes using a tmpl
	ts, err := template.ParseFiles("./ui/html/show_page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		log.Print(q[0])
		return
	}
}

func (app *application) createQuoteForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/quotes_form_page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createQuote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	author := r.PostForm.Get("author_name")
	category := r.PostForm.Get("category")
	quote := r.PostForm.Get("quote")
	//check the web form fields to ladidate
	errors := make(map[string]string)
	//check each field
	if strings.TrimSpace(author) == "" {
		errors["author"] = "Thsi field cannot be left blank"
	} else if utf8.RuneCountInString(author) > 50 {
		errors["author"] = "This field is too long, max lenght is 50 characters"
	}
	if strings.TrimSpace(category) == "" {
		errors["category"] = "This filed cannot be left blank"
	} else if utf8.RuneCountInString(category) > 25 {
		errors["category"] = "This field is too long, max lenght is 25 characters"
	}
	if strings.TrimSpace(quote) == "" {
		errors["quote"] = "This filed cannot be left blank"
	} else if utf8.RuneCountInString(category) > 255 {
		errors["quote"] = "This field is too long, max lenght is 255 characters"
	}

	//check if there are any errors in the map

	if len(errors) > 0 {
		ts, err := template.ParseFiles("./ui/html/quotes_form_page.tmpl")
		if err != nil {
			log.Println(err.Error())
			app.clientError(w, http.StatusBadRequest)
			return
		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm: errors,
			FormData:       r.PostForm,
		})
		if err != nil {
			app.serverError(w, err)
			return
		}
		return
	}

	//Insert a Quote
	id, err := app.quotes.Insert(author, category, quote)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/quote/%d", id), http.StatusSeeOther)
}

func (app *application) showQuote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	q, err := app.quotes.Getid(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//an instance of templeData
	data := &templateData{
		Quote: q,
	}
	// Display the quote using a template
	ts, err := template.ParseFiles("./ui/html/quote_page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}

}
