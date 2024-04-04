package main

import (
	"net/url"

	"amencia.net/quotebox/pkg/models"
)

type templateData struct {
	Quotes         []*models.Quote
	ErrorsFromForm map[string]string
	FormData       url.Values
}
