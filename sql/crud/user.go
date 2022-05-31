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

func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	var user *models.User
	db.First(&user, userID)
	if user == nil {
		return nil, fmt.Errorf("未找到用户")
	}
	return user, nil
}

func GetUserByName(db *gorm.DB, name string) *models.User {
	var user *models.User
	db.Where(&models.User{Name: name}).First(&user)
	return user
}

func SubscribeUser(db *gorm.DB, userID uint, subscriberUserID uint) (*models.User, error) {
	var subscriber, user *models.User
	db.First(&subscriber, subscriberUserID)
	db.First(&user, userID)
	if subscriber == nil {
		return nil, fmt.Errorf("未找到关注人")
	}
	if user == nil {
		return nil, fmt.Errorf("未找到用户")
	}
	db.Model(&user).Association("Subscribers").Append(subscriber)
	return user, nil
}

func CancelSubscribeUser(db *gorm.DB, userID uint, subscriberUserID uint) (*models.User, error) {
	var subscriber, user *models.User
	db.First(&subscriber, subscriberUserID)
	db.First(&user, userID)
	if subscriber == nil {
		return nil, fmt.Errorf("未找到关注人")
	}
	if user == nil {
		return nil, fmt.Errorf("未找到用户")
	}
	if err := db.Model(&user).Association("Subscribers").Delete(subscriber); err != nil {
		return nil, fmt.Errorf("关注不存在")
	}
	return user, nil
}

func GetUserSubscribersByID(db *gorm.DB, userID uint) []models.User {
	var user *models.User
	db.Preload("Subscribers").Find(&user, userID)
	return user.Subscribers
}

func GetUserSubscribersCountByID(db *gorm.DB, userID uint) int64 {
	var count int64
	db.Raw("select count(subscriber_id) from subscribes where user_id = ?", &userID).Scan(&count)
	return count
}

func GetUserFollowersByID(db *gorm.DB, userID uint) []models.User {
	var followers []models.User
	db.Raw("select * from users where id in"+
		"(select user_id from subscribes left join `users`"+
		"on `users`.id = subscriber_id "+
		"where subscriber_id = ?)", userID).Scan(&followers)
	return followers
}

func GetUserFollowersCountByID(db *gorm.DB, userID uint) int64 {
	var count int64
	db.Raw("select count(user_id) from subscribes where subscriber_id = ?", &userID).Scan(&count)
	return count
}

func IsUserFollow(db *gorm.DB, userID, anotherUserID uint) bool {
	if userID == anotherUserID {
		return false
	}
	var user *models.User
	db.Raw("select * from users where id in"+
		" (select user_id from subscribes where user_id = ? and subscriber_id = ?)", userID, anotherUserID).Scan(&user)
	return user != nil
}
