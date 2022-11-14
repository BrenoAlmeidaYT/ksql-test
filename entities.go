package main

import "github.com/vingarcia/ksql"

var VideosTable = ksql.NewTable("videos", "id")
var AuthorsTable = ksql.NewTable("authors", "id")

type Video struct {
	ID          int     `ksql:"id" json:"id"`
	Title       *string `ksql:"title" json:"title"`
	Description *string `ksql:"description" json:"description"`
	LikeCount   *int    `ksql:"like_count" json:"like_count"`
	ViewCount   *int    `ksql:"view_count" json:"view_count"`
	AuthorID    *int    `ksql:"author_id" json:"author_id"`
}

type Author struct {
	ID    int     `ksql:"id" json:"id"`
	Name  *string `ksql:"name" json:"name"`
	Phone *string `ksql:"phone" json:"phone"`
}
