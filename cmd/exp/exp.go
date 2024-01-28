package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	id      int
	title   string
	content string
	author  string
}

var conn *pgx.Conn

func main() {
	var err error
	dbURL := "postgres://postgres:secret@localhost:5432/postgres"
	conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conexão com o banco foi efetuada com sucesso")
	defer conn.Close(context.Background())

	// createTable()
	// insertPost()
	// insertPostWithReturn()
	// selectPostById()
	selectAllPosts()
}

func createTable() {
	query := `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT, 
			author TEXT NOT NULL
		);
	`
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table posts created")
}

// SQL INJECTION
func insertPost() {
	title := "','',''); DROP TABLE posts; --"
	content := "Conteúdo do post 1"
	author := "robson"
	// query := fmt.Sprintf(`
	// 	INSERT INTO posts (title, content, author)
	// 	values ('%s', '%s', '%s')
	// `, title, content, author)

	query := `
		INSERT INTO posts (title, content, author)
		values ($1, $2, $3)
	`
	//fmt.Println(query)
	_, err := conn.Exec(context.Background(), query, title, content, author)
	if err != nil {
		panic(err)
	}
	fmt.Println("Post criado com sucesso")
}

func insertPostWithReturn() {
	title := "Post 3"
	content := "Conteúdo do post 3"
	author := "robson"
	query := `
		INSERT INTO posts (title, content, author)
		values ($1, $2, $3) RETURNING id;
	`
	row := conn.QueryRow(context.Background(), query, title, content, author)
	var id int
	if err := row.Scan(&id); err != nil {
		panic(err)
	}
	fmt.Println("Post criado. Id =", id)
}

func selectPostById() {
	id := 4
	var title, content, author string
	query := "select title, content, author from posts where id = $1"
	row := conn.QueryRow(context.Background(), query, id)
	err := row.Scan(&title, &content, &author)
	if err == pgx.ErrNoRows {
		fmt.Println("No post found for id =", id)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("Post: title=%s, content=%s, author=%s \n", title, content, author)
}

func selectAllPosts() {
	query := "select id, title, content, author from posts"
	// query := "select * from posts"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Err() != nil {
		panic(rows.Err())
	}

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.id, &post.title, &post.content, &post.author)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	for _, post := range posts {
		fmt.Println(post)
	}
}
