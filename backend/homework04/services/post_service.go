package services

import (
	"homework04/models"
	"homework04/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}
func (p *PostService) GetPosts() ([]gin.H, error) {
	var posts []models.Post
	if err := p.db.Preload("User").Find(&posts).Error; err != nil {
		return nil, err
	}

	data := make([]gin.H, len(posts))
	for i, p := range posts {
		data[i] = gin.H{
			"id":      p.ID,
			"title":   p.Title,
			"content": p.Content,
			"author":  p.User.Username,
			"created": p.CreatedAt,
		}
	}
	return data, nil

}
func (p *PostService) CreatePost(post models.Post) (uint, error) {
	if err := p.db.Create(&post).Error; err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (p *PostService) GetPost(id string) (models.Post, error) {
	var post models.Post
	if err := p.db.Preload("User").First(&post, id).Error; err != nil {
		return post, utils.NewBizError(500, "post not found")
	}
	return post, nil
}
func (p *PostService) UpdatePost(userID uint, id string, req models.UpdatePostRequest) error {
	var post models.Post
	if err := p.db.First(&post, id).Error; err != nil {
		return utils.NewBizError(404, "post not found")
	}

	if post.UserID != userID {
		return utils.NewBizError(404, "post not found")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := p.db.Save(&post).Error; err != nil {
		return utils.NewBizError(500, "failed to update post")
	}
	return nil
}

func (p *PostService) DeletePost(userID uint, id string) error {
	var post models.Post
	if err := p.db.First(&post, id).Error; err != nil {
		return utils.NewBizError(404, "post not found")
	}

	if post.UserID != userID {
		return utils.NewBizError(403, "forbidden: not the author")
	}

	if err := p.db.Delete(&post).Error; err != nil {
		return utils.NewBizError(500, "failed to delete post")
	}
	return nil
}
