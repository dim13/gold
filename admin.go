package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dim13/gold/articles"
	"github.com/dim13/gold/storage"
)

type adminPage struct {
	Articles articles.Articles
	Article  *articles.Article
	Title    string
	Config   storage.Config
	Error    string
}

func (p adminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Config = conf
	err := tmpl.ExecuteTemplate(w, "admin.tmpl", p)
	if err != nil {
		log.Println(err)
	}
}

type adminIndex struct{ adminPage }

func (p *adminIndex) Select(_ []string) bool {
	p.Articles = art
	p.Title = "Admin Interface"
	return true
}

func (p *adminIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p.adminPage.ServeHTTP(w, r)
		return
	}
	switch r.FormValue("submit") {
	case "add":
		title := r.FormValue("title")
		r.URL.Path = "/admin/" + articles.MakeSlug(title)
	}
}

type adminSlug struct{ adminPage }

func (p *adminSlug) Select(match []string) bool {
	if a, ok := art.Find(match[0]); ok {
		p.Title = a.Title
		p.Article = a
	} else {
		p.Article = &articles.Article{
			Slug:  match[0],
			Title: articles.MakeTitle(match[0]),
			Date:  time.Now(),
		}
	}
	return true
}

func (p *adminSlug) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p.adminPage.ServeHTTP(w, r)
		return
	}
	a := articles.Article{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Tags:    articles.ReadTags(r.FormValue("tags")),
		Body:    r.FormValue("body"),
		Enabled: r.FormValue("enabled") == "on",
	}
	if p.Article.Slug != a.Slug {
		art.Delete(*p.Article)
		r.URL.Path = "/admin/" + a.Slug
	}
	switch r.FormValue("submit") {
	case "reload":
		log.Println("reloading")
		art.Load()
	case "preview":
		art.Add(a)
	case "save":
		art.Add(a)
		art.Store()
		r.URL.Path = "/admin/"
	case "delete":
		art.Delete(a)
		art.Store()
		r.URL.Path = "/admin/"
	case "cancel":
		r.URL.Path = "/admin/"
	}
}
