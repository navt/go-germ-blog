package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./packs/models"
	"./packs/utility"
)

var posts map[string]*models.Post

func indexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/index.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "index", posts)
}

func writePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/write.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "write", nil)
}

func editPage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	onePost, found := posts[id]
	if found == false {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("tmpl/write.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "write", onePost)
}

func showPage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	onePost, found := posts[id]
	if found == false {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("tmpl/show.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "show", onePost)
}

func remove(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	_, found := posts[id]
	if found == false {
		http.NotFound(w, r)
		return
	}
	delete(posts, id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.EscapedPath()
	query = query[1:]
	switch query {
	case "index":
		indexPage(w, r)
		return
	case "":
		indexPage(w, r)
		return
	case "show":
		showPage(w, r)
		return
	case "write":
		id := r.FormValue("id")
		if id == "" {
			writePage(w, r)
		} else {
			editPage(w, r)
		}
		return
	case "delete":
		remove(w, r)
		return
	}
	// если запрос не определён
	http.NotFound(w, r)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.EscapedPath()
	query = query[1:]
	if query == "save" {
		id := r.FormValue("id")
		title := r.FormValue("title")
		description := r.FormValue("description")
		content := r.FormValue("content")
		var onePost *models.Post
		// если ещё нет такого поста
		if id == "" {
			id = utility.IdFromTime()
			onePost = models.CreatePost(id, title, description, content)
			posts[onePost.Id] = onePost
		} else {
			onePost = posts[id]
			// перезапись полей
			onePost.Title = title
			onePost.Description = description
			onePost.Content = content
		}
		http.Redirect(w, r, "/", http.StatusFound) // 302
	}
	// если это не save
	http.NotFound(w, r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		PostHandler(w, r)
		return
	}
	http.NotFound(w, r)
}

func main() {
	fmt.Println("Слушаем порт :8000")
	posts = make(map[string]*models.Post, 0)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8000", nil)
}
