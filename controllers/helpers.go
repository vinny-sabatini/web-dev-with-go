package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	dec := schema.NewDecoder()
	err := dec.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}
	return nil
}
