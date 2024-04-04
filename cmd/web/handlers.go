package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
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
	w.Write([]byte("Welcome to quotebox."))
}

func (app *application) createQuoteForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/quotes_form_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) createQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/quote", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
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
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm: errors,
			FormData:       r.PostForm,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		return
	}

	//Insert a Quote
	id, err := app.quotes.Insert(author, category, quote)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
	}
	fmt.Print(w, "row with id: %d has been inserted.", id)

}

func (app *application) displayQuotation(w http.ResponseWriter, r *http.Request) {
	q, err := app.quotes.Read()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	data := &templateData{
		Quotes: q,
	}

	//Display Quotes using a tmpl
	ts, err := template.ParseFiles("./ui/html/show_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		log.Print(q[0])
		return
	}
}
