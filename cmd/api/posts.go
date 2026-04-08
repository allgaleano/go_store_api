package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/allgaleano/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title 	string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=100"`
	Tags 	[]string `json:"tags" validate:"required,max=100"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := &store.Post{
		Title: payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
			case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			default:
			app.internalServerError(w, r, err)
		}
		return
	}
	
	comments, err := app.store.Comments.GetByPostID(ctx, id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments
	
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title 	*string  `json:"title" validate:"omitempty,max=100"`
	Content *string  `json:"content" validate:"omitempty,max=100"`
	Tags 	[]string `json:"tags" validate:"omitempty,max=100"`
}

func (app *application) patchPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := app.store.Posts.PatchByID(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	
	ctx := r.Context()
	
	if err := app.store.Posts.DeleteByID(ctx, id); err != nil {
		switch {
			case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
