package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Category struct {
	Id     uint    `json:"id"`
	Name   string  `json:"name"`
	Groups []Group `json:"groups"`
}

type Group struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Posts []Post `json:"posts"`
}

type Post struct {
	Id       uint      `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	UID   uint   `json:"uid"`
}

func (a *App) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query(context.Background(), "select name, id from categories")

	defer rows.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while fetching the data (categories): %v", err)
		return
	}

	var categories []Category

	for rows.Next() {
		var id uint
		var name string

		err := rows.Scan(&name, &id)

		if err != nil {
			fmt.Printf("Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error while fetching the data (Categories): %v", err)
			return
		}

		categories = append(categories, Category{
			Id:   id,
			Name: name,
		})
	}

	for i := 0; i < len(categories); i++ {
		g, err := a.GetGroups(categories[i].Id)

		if err != nil {
			fmt.Printf("Piz gets %v", err)
		}

		categories[i].Groups = g
	}

	j, _ := json.Marshal(categories)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "%v", string(j))
}

func (a *App) GetGroups(id uint) (groups []Group, e error) {
	rows, err := a.DB.Query(context.Background(), "select title, id as gid from groups where category_id = $1", &id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups = []Group{}

	for rows.Next() {
		var title string
		var gid uint

		er := rows.Scan(&title, &gid)

		if er != nil {
			return nil, er
		}

		groups = append(groups, Group{
			Id:    gid,
			Title: title,
		})
	}

	for i := 0; i < len(groups); i++ {
		p, err := a.GetPosts(groups[i].Id)

		if err != nil {
			fmt.Printf("Error getting comment %v", err)
		}

		groups[i].Posts = p
	}

	return groups, nil
}

func (a *App) GetPosts(id uint) (posts []Post, e error) {
	rows, err := a.DB.Query(context.Background(), "select title, content, id from posts where group_id = $1", &id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts = []Post{}

	for rows.Next() {
		var title string
		var content string
		var pid uint

		er := rows.Scan(&title, &content, &pid)

		if er != nil {
			return nil, er
		}

		posts = append(posts, Post{
			Id:      pid,
			Title:   title,
			Content: content,
		})
	}

	for i := 0; i < len(posts); i++ {
		c, err := a.GetComments(posts[i].Id)

		if err != nil {
			fmt.Printf("Error getting comment %v", err)
		}

		posts[i].Comments = c
	}

	return posts, nil
}

func (a *App) GetComments(id uint) ([]Comment, error) {
	r, err := a.DB.Query(context.Background(),
		`select
		id as cid,
		title,
		user_id as uid
	from 
		comments 
	where 
		post_id = $1`, &id)

	defer r.Close()

	var comments []Comment

	if err != nil {
		fmt.Printf("Mo rows")
	}

	for r.Next() {
		var title string
		var cid uint
		var uid uint

		er := r.Scan(&cid, &title, &uid)

		if er != nil {
			return nil, er
		}

		comments = append(comments, Comment{
			Id:    cid,
			Title: title,
			UID:   uid,
		})
	}

	if comments == nil {
		comments = []Comment{}
	}

	return comments, nil
}
