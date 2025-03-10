package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanisharrma/gopher-social/internal/store"
)

type postKey string

var postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title"   validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// GetPost
//
//	@Summary		Fetches Post details
//	@Description	Getches Post details by post ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int			true	"Post ID"
//	@Success		200	{object}	store.Post	"Post fetched"
//	@Failure		400	{object}	error		"Payload missing"
//	@Failure		401	{object}	error		"Unauthorized"
//	@Failure		404	{object}	error		"Not Found"
//	@Failure		500	{object}	error		"An error occured
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, "Post fetched", post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Create a post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post Payload"
//	@Success		201		{object}	store.Post			"Post created"
//	@Failure		400		{object}	error				"Payload missing"
//	@failure		401		{object}	error				"Unauthorized"
//	@Failure		500		{object}	error				"An error occured"
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromCtx(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "Post created successfully", post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Deletes a post using post ID by Authorized (admin, owner)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Post ID"
//	@Success		204	{string}	string	"Post deleted"
//	@Failure		400	{object}	error	"Payload missing"
//	@Failure		401	{object}	error	"Unauthorized"
//	@Failure		404	{object}	error	"Not Found"
//	@Failure		500	{obejct}	error	"An error occured"
//	@Security		ApiKeyAuth
//	@router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid params"))
	}

	if err := app.store.Posts.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return

	}

	if err := app.jsonResponse(w, http.StatusNoContent, "Post deleted succesfully", nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title   *string   `json:"title"   validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags"    validate:"omitempty"`
}

// UpdatePost
//
//	@Summary		Update a post
//	@Description	Updates a post using post ID by authorized (admin,moderator,owner)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"UpdatePost payload"
//	@Success		200		{string}	string				"Post updated"
//	@Failure		400		{object}	error				"Payload missing"
//	@Failure		401		{object}	error				"Unauthorized"
//	@Failure		404		{object}	error				"Not Found"
//	@Failure		500		{object}	error				"An error occured"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Println("update payload", payload)
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	err := app.store.Posts.Update(r.Context(), post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "Post updated successfully", nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid params"))
			return
		}
		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
