package handlers

import (
	"homework04/models"
	"homework04/services"
	"homework04/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (p *PostHandler) GetPosts(c *gin.Context) {
	data, err := p.postService.GetPosts()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "ok",
		Data:    data,
	})
}
func (p *PostHandler) CreatePost(c *gin.Context) {
	// 创建文章
	userID := c.GetUint("user_id")
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBizError(400, err.Error()))
		return
	}
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	id, err := p.postService.CreatePost(post)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "post created",
		Data:    gin.H{"post_id": id},
	})
}

// 获取单篇文章
func (p *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := p.postService.GetPost(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "ok",
		Data: gin.H{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"author":  post.User.Username,
			"created": post.CreatedAt,
		},
	})
}

// 更新文章
func (p *PostHandler) UpdatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBizError(400, err.Error()))
		return
	}
	err := p.postService.UpdatePost(userID, id, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "post updated",
	})
}

// 删除文章
func (p *PostHandler) DeletePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	err := p.postService.DeletePost(userID, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "post deleted",
	})
}
