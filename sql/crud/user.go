package crud

import (
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *sql.User) *sql.User {
	db.Create(&user).Commit()
	return user
}

func GetUser(db *gorm.DB, id uint64) *sql.User {
	var user *sql.User
	db.First(&user, id)
	return user
}

func GetUserByName(db *gorm.DB, name string) *sql.User {
	var user *sql.User
	db.Where(&sql.User{Name: name}).First(&user)
	return user
}

func GetUserSubscribersByName(db *gorm.DB, name string) *sql.User {
	var user *sql.User
	db.Where(&sql.User{Name: name}).Preload("Subscribers").Find(&user)
	return user
}

func GetUserFollowersByName(db *gorm.DB, name string) *sql.User {
	user := GetUserByName(db, name)
	db.Raw("select * from users where id in"+
		"(select user_id from subscribes, `users`"+
		"where `users`.id = subscriber_id and subscriber_id = ?)", user.ID).Scan(&user.Followers)
	return user
}
