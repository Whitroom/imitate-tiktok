package crud

import (
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) *models.User {
	db.Create(&user).Commit()
	return user
}

func GetUser(db *gorm.DB, id uint64) *models.User {
	var user *models.User
	db.First(&user, id)
	return user
}

func GetUserByName(db *gorm.DB, name string) *models.User {
	var user *models.User
	db.Where(&models.User{Name: name}).First(&user)
	return user
}

func GetUserSubscribersByName(db *gorm.DB, name string) *models.User {
	var user *models.User
	db.Where(&models.User{Name: name}).Preload("Subscribers").Find(&user)
	return user
}

func GetUserFollowersByName(db *gorm.DB, name string) *models.User {
	user := GetUserByName(db, name)
	db.Raw("select * from users where id in"+
		"(select user_id from subscribes, `users`"+
		"where `users`.id = subscriber_id and subscriber_id = ?)", user.ID).Scan(&user.Followers)
	return user
}
