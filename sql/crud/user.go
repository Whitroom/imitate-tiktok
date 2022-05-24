package crud

import (
	"fmt"

	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) *models.User {
	db.Create(&user).Commit()
	return user
}

func GetUserByID(db *gorm.DB, user_id uint64) (*models.User, error) {
	var user *models.User
	db.First(&user, user_id)
	if user == nil {
		err := fmt.Errorf("未找到用户")
		return nil, err
	}
	return user, nil
}

func GetUsersByName(db *gorm.DB, name string) []models.User {
	var users []models.User
	db.Where(&models.User{Name: name}).First(&users)
	return users
}

func SubscribeUser(db *gorm.DB, user_id uint, subscriber_user_id uint) (*models.User, error) {
	var subscriber, user *models.User
	db.First(&subscriber, subscriber_user_id)
	db.First(&user, user_id)
	if subscriber == nil {
		err := fmt.Errorf("未找到关注人")
		return nil, err
	}
	if user == nil {
		err := fmt.Errorf("未找到用户")
		return nil, err
	}
	db.Model(&user).Association("Subscriber").Append(&subscriber)
	return user, nil
}

func GetUserSubscribersByID(db *gorm.DB, user_id uint) []models.User {
	var user *models.User
	db.Preload("Subscribers").Find(&user, user_id)
	return user.Subscribers
}

func GetUserFollowersByName(db *gorm.DB, user_id uint) []models.User {
	var followers []models.User
	db.Raw("select * from users where id in"+
		"(select user_id from subscribes left join `users`"+
		"on `users`.id = subscriber_id "+
		"where subscriber_id = ?)", user_id).Scan(&followers)
	return followers
}
