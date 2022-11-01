package main

import (
	"fmt"
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
