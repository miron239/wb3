package cockroach

import (
	"go-microservices/libs/logger"
	"time"
)

type Post struct {
	ID            int64 `gorm:"primary_key"`
	UserId        string
	PostpaidLimit int
	Spp           int
	ShippingFee   int
	ReturnFee     int
	CreatedAt     time.Time
}

func GetAllPosts() ([]Post, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}

	var posts []Post
	err = db.Find(&posts).Error
	if err != nil {
		logger.GetCockroach().Error(11, map[string]interface{}{
			"error":  err,
			"method": "GetAllPosts",
		})
		return nil, err
	}

	return posts, nil
}

func GetPost(id int64) (*Post, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}

	post := Post{}
	err = db.Where(&Post{ID: id}).First(&post).Error
	if err != nil {
		logger.GetCockroach().Error(11, map[string]interface{}{
			"error":  err,
			"method": "GetPost",
			"id":     id,
		})
		return nil, err
	}

	return &post, nil
}

func CreatePost(p Post) (*Post, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}

	post := Post{}
	err = db.Create(&Post{
		CreatedAt:     time.Now(),
		ID:            p.ID,
		PostpaidLimit: p.PostpaidLimit,
		Spp:           p.Spp,
		ShippingFee:   p.ShippingFee,
		ReturnFee:     p.ReturnFee,
	}).Scan(&post).Error
	if err != nil {
		logger.GetCockroach().Error(11, map[string]interface{}{
			"error":  err,
			"method": "CreatePost",
			"post":   p,
		})
		return nil, err
	}

	return &post, nil
}

func UpdatePost(q Post, u Post) error {
	db, err := Connection()
	if err != nil {
		return err
	}

	err = db.Model(&q).Updates(&Post{
		ID:            u.ID,
		PostpaidLimit: u.PostpaidLimit,
		Spp:           u.Spp,
		ShippingFee:   u.ShippingFee,
		ReturnFee:     u.ReturnFee,
		CreatedAt:     time.Now(),
	}).Error
	if err != nil {
		logger.GetCockroach().Error(11, map[string]interface{}{
			"error":  err,
			"method": "UpdatePost",
			"query":  q,
			"update": u,
		})
		return err
	}

	return nil
}

func DeletePost(id int64) error {
	db, err := Connection()
	if err != nil {
		return err
	}

	err = db.Delete(&Post{ID: id}).Error
	if err != nil {
		logger.GetCockroach().Error(11, map[string]interface{}{
			"error":  err,
			"method": "DeletePost",
			"id":     id,
		})
		return err
	}

	return nil
}
