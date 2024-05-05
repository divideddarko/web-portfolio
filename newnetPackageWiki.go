// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"log"
	"net/http"
	"os"
	"html/template"
	"regexp"
	"fmt"
)

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html", "tmpl/frontPage.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}


func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/*func rootHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err:= loadPage(title)
  fmt.Println("helo", err)
  if err != nil {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }
  renderTemplate(w, "frontPage", p)
}
*/

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    fmt.Println(m, "mmmm")
    if m == nil {
      /*http.NotFound(w, r)*/
      body := r.FormValue("body")
      p := &Page{Title: "hello", Body: []byte(body)}
      if p == nil {
	fmt.Println("p nil", p)
      }

      if p != nil {
	fmt.Println("p pop", p)
      }


      http.Redirect(w, r, "/", http.StatusFound)
      fmt.Println("w: ", w, "r: ", r)

      renderTemplate(w, "frontPage", p)
      /*return*/
    }
    fn(w, r, m[2])
  }
}

func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  log.Fatal(http.ListenAndServe(":8089", nil))
}
