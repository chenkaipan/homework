package services

import (
	"homework04/models"
	"homework04/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}
func (cs *CommentService) CreateComment(userID uint, postID string, req models.CreateCommentRequest) (uint, error) {
	var post models.Post
	if err := cs.db.First(&post, postID).Error; err != nil {
		return 0, utils.NewBizError(404, "post not found")
	}
	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  post.ID,
	}

	if err := cs.db.Create(&comment).Error; err != nil {
		return 0, utils.NewBizError(500, "failed to create comment")
	}
	return comment.ID, nil
}

func (cs *CommentService) ListComments(postID string) ([]gin.H, error) {
	var comments []models.Comment
	if err := cs.db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, utils.NewBizError(500, "failed to list comments")
	}

	data := make([]gin.H, len(comments))
	for i, cm := range comments {
		data[i] = gin.H{
			"id":      cm.ID,
			"content": cm.Content,
			"author":  cm.User.Username,
			"created": cm.CreatedAt,
		}
	}
	return data, nil

}
