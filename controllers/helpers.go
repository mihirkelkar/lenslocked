package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

func ParseForm(r *http.Request, dst interface{}) error {
	//Note that the dst parameter that should be passed in here.
	//should be a pointer to a struct and not a struct itself
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	decoder := schema.NewDecoder()
	if err := decoder.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
