package models

type Post struct {
	Id          string
	Title       string
	Description string
	Content     string
}

func CreatePost(id, title, description, content string) *Post {
	return &Post{id, title, description, content}
}
