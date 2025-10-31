package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
	"github.com/yakuter/ugin/internal/service"
)

type PostHandler struct {
	service service.PostService
}

// NewPostHandler creates a new post handler
func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

// GetByID handles GET /posts/:id
// @Summary Get post by ID
// @Description Get a single post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} domain.Post
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	post, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// List handles GET /posts
// @Summary List posts
// @Description Get all posts with pagination and filtering
// @Tags posts
// @Accept json
// @Produce json
// @Param Limit query int false "Limit" default(25)
// @Param Offset query int false "Offset" default(0)
// @Param Sort query string false "Sort field" default(id)
// @Param Order query string false "Sort order" default(DESC)
// @Param Search query string false "Search keyword"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/posts [get]
func (h *PostHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("Limit", "25"))
	offset, _ := strconv.Atoi(c.DefaultQuery("Offset", "0"))

	filter := repository.ListFilter{
		Search: c.Query("Search"),
		Limit:  limit,
		Offset: offset,
		Sort:   c.DefaultQuery("Sort", "id"),
		Order:  c.DefaultQuery("Order", "DESC"),
	}

	posts, result, err := h.service.List(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          posts,
		"total_data":    result.Total,
		"filtered_data": result.Filtered,
	})
}

// Create handles POST /posts
// @Summary Create post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body domain.CreatePostRequest true "Post object"
// @Success 201 {object} domain.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var post domain.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := h.service.Create(ctx, &post); err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// Update handles PUT /posts/:id
// @Summary Update post
// @Description Update an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body domain.CreatePostRequest true "Post object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var post domain.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if err := h.service.Update(ctx, id, &post); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if errors.Is(err, repository.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
}

// Delete handles DELETE /posts/:id
// @Summary Delete post
// @Description Delete a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if err := h.service.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully", "id": id})
}

