package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

func main() {
	connStr := "user=postgres password=password dbname=ksqldb sslmode=disable host=localhost port=8999"
	ctx := context.Background()
	db, err := kpgx.New(ctx, connStr, ksql.Config{
		MaxOpenConns: 100,
	})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	ctrl := NewController(db)

	app.Post("/authors", ctrl.InsertAuthor)
	app.Post("/authors/:author_id/videos", ctrl.InsertVideo)

	app.Patch("/authors", ctrl.UpdateAuthor)
	app.Patch("/videos", ctrl.UpdateVideo)

	app.Delete("/authors/:author_id", ctrl.DeleteAuthor)
	app.Delete("/videos/:video_id", ctrl.DeleteVideo)

	app.Get("/authors/:author_id", ctrl.GetAuthor)

	app.Get("/authors/phone/test/chunks", ctrl.GetAuthorWithInvalidPhonesChunks)
	app.Get("/authors/phone/test/nochunks", ctrl.GetAuthorWithInvalidPhonesWithoutChunks)

	app.Post("/authorWithVideos", ctrl.InsertAuthorWithVideos)

	app.Listen(":3100")
}
