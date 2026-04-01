package handlers

import (
	"homework04/models"
	"homework04/services"
	"homework04/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// 创建评论
func (ch *CommentHandler) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID := c.Param("post_id")
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBizError(400, err.Error()))
		return
	}

	id, err := ch.commentService.CreateComment(userID, postID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "comment created",
		Data:    gin.H{"comment_id": id},
	})
}

// 获取某篇文章的评论列表
func (ch *CommentHandler) ListComments(c *gin.Context) {
	postID := c.Param("post_id")

	data, err := ch.commentService.ListComments(postID)
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
