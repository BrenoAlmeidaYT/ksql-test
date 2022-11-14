package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ksql"
)

type Controller struct {
	db ksql.DB
}

func NewController(db ksql.DB) Controller {
	return Controller{
		db: db,
	}
}

func (c Controller) InsertAuthor(ctx *fiber.Ctx) error {
	var author Author

	if err := ctx.BodyParser(&author); err != nil {
		return err
	}

	err := c.db.Insert(ctx.Context(), AuthorsTable, &author)
	if err != nil {
		return nil
	}
	return ctx.JSON(author)
}

func (c Controller) InsertVideo(ctx *fiber.Ctx) error {
	var video Video

	if err := ctx.BodyParser(&video); err != nil {
		return nil
	}
	authorID, err := strconv.Atoi(ctx.Params("author_id"))
	if err != nil {
		return err
	}
	video.AuthorID = &authorID

	err = c.db.Insert(ctx.Context(), VideosTable, &video)
	if err != nil {
		return err
	}

	return ctx.JSON(video)
}

func (c Controller) UpdateAuthor(ctx *fiber.Ctx) error {
	var author Author

	if err := ctx.BodyParser(&author); err != nil {
		return err
	}

	err := c.db.Patch(ctx.Context(), AuthorsTable, &author)
	if err != nil {
		if err == ksql.ErrRecordNotFound {
			return ctx.Status(fiber.ErrNotFound.Code).SendString(fmt.Sprintf(`Not found author with the ID "%d"`, author.ID))
		}
		return err
	}

	return ctx.JSON(author)
}

func (c Controller) UpdateVideo(ctx *fiber.Ctx) error {
	var video Video

	if err := ctx.BodyParser(&video); err != nil {
		return err
	}

	err := c.db.Patch(ctx.Context(), VideosTable, &video)
	if err != nil {
		if err == ksql.ErrRecordNotFound {
			return ctx.Status(fiber.ErrNotFound.Code).SendString(fmt.Sprintf(`Not found video with the ID "%d"`, video.ID))
		}
		return err
	}

	return ctx.JSON(video)
}

func (c Controller) DeleteAuthor(ctx *fiber.Ctx) error {
	err := c.db.Delete(ctx.Context(), AuthorsTable, ctx.Params("author_id"))
	if err != nil {
		if err == ksql.ErrRecordNotFound {
			return ctx.Status(fiber.ErrNotFound.Code).SendString(fmt.Sprintf(`Not found author with the ID "%s"`, ctx.Params("author_id")))
		}
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (c Controller) DeleteVideo(ctx *fiber.Ctx) error {
	err := c.db.Delete(ctx.Context(), VideosTable, ctx.Params("video_id"))
	if err != nil {
		if err == ksql.ErrRecordNotFound {
			return ctx.Status(fiber.ErrNotFound.Code).SendString(fmt.Sprintf(`Not found video with the ID "%s"`, ctx.Params("video_id")))
		}
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (c Controller) GetAuthor(ctx *fiber.Ctx) error {
	var author Author
	err := c.db.QueryOne(ctx.Context(), &author, "FROM authors WHERE ID = $1", ctx.Params("author_id"))
	if err != nil {
		if err == ksql.ErrRecordNotFound {
			return ctx.Status(fiber.ErrNotFound.Code).SendString(fmt.Sprintf(`Not found author with the ID "%s"`, ctx.Params("author_id")))
		}
		return err
	}
	return ctx.JSON(author)
}

func (c Controller) GetAuthorWithInvalidPhonesWithoutChunks(ctx *fiber.Ctx) error {
	fmt.Println("inicio sem chunks")
	PrintMemoryUsage()
	var authors []Author = []Author{}
	err := c.db.Query(ctx.Context(), &authors, "FROM authors")
	if err != nil {
		return err
	}

	var authorsWithInvalidPhones []Author = []Author{}

	for _, author := range authors {
		if len(*author.Phone) != 11 {
			authorsWithInvalidPhones = append(authorsWithInvalidPhones, author)
		}
	}
	PrintMemoryUsage()
	fmt.Println("final sem chunks")
	return ctx.JSON(authorsWithInvalidPhones)
}

func (c Controller) GetAuthorWithInvalidPhonesChunks(ctx *fiber.Ctx) error {
	fmt.Println("inicio com chunks")
	PrintMemoryUsage()

	var authorsWithInvalidPhones []Author = []Author{}
	err := c.db.QueryChunks(ctx.Context(), ksql.ChunkParser{
		Query:     "FROM authors",
		Params:    []interface{}{},
		ChunkSize: 100,
		ForEachChunk: func(authors []Author) error {
			for _, author := range authors {
				if len(*author.Phone) != 11 {
					authorsWithInvalidPhones = append(authorsWithInvalidPhones, author)
				}
			}
			return nil
		},
	})
	if err != nil {
		return err
	}

	PrintMemoryUsage()
	fmt.Println("final com chunks")
	return ctx.JSON(authorsWithInvalidPhones)
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB\n", m.Sys/1024/1024)
}

func (c Controller) InsertAuthorWithVideos(ctx *fiber.Ctx) error {
	var input struct {
		Name   *string `json:"author_name"`
		Phone  *string `json:"author_phone"`
		Videos []struct {
			Title       *string `json:"title"`
			Description *string `json:"description"`
			LikeCount   *int    `json:"like_count"`
			ViewCount   *int    `json:"view_count"`
		} `json:"videos"`
	}
	var output struct {
		AuthorID int     `json:"author_id"`
		Name     *string `json:"author_name"`
		Phone    *string `json:"author_phone"`
		Videos   []Video `json:"videos"`
	}

	var videos []Video

	if err := ctx.BodyParser(&input); err != nil {
		return err
	}

	err := c.db.Transaction(ctx.Context(), func(db ksql.Provider) error {
		var author Author = Author{
			Name:  input.Name,
			Phone: input.Phone,
		}
		err := db.Insert(ctx.Context(), AuthorsTable, &author)
		if err != nil {
			return err
		}
		output.AuthorID = author.ID
		output.Name = author.Name
		output.Phone = author.Phone

		for _, v := range input.Videos {
			video := Video{
				Title:       v.Title,
				Description: v.Description,
				LikeCount:   v.LikeCount,
				ViewCount:   v.ViewCount,
				AuthorID:    &author.ID,
			}
			err := db.Insert(ctx.Context(), VideosTable, &video)
			if err != nil {
				return err
			}
			videos = append(videos, video)
		}

		// aqui a transação é commitada
		return nil
	})
	if err != nil {
		return err
	}
	output.Videos = videos

	return ctx.JSON(output)
}
