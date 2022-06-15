package crud

import (
	"fmt"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
)

func CreateUser(user *models.User) *models.User {
	statement := sql.DB.Create(&user)
	if err := statement.Error; err != nil {
		return nil
	}
	statement.Commit()
	return user
}

func GetUserByID(userID uint) (*models.User, error) {
	var user *models.User
	sql.DB.First(&user, userID)
	if user == nil {
		return nil, fmt.Errorf("未找到用户")
	}
	return user, nil
}

func GetUserByName(name string) *models.User {
	var user *models.User
	err := sql.DB.Where(&models.User{Name: name}).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return user
}

func SubscribeUser(userID uint, subscriberUserID uint) (*models.User, error) {
	var subscriber, user *models.User
	sql.DB.First(&subscriber, subscriberUserID)
	sql.DB.First(&user, userID)
	if subscriber == nil {
		return nil, fmt.Errorf("未找到关注人")
	}
	if user == nil {
		return nil, fmt.Errorf("未找到用户")
	}
	if err := sql.DB.Model(&user).Association("Subscribers").Append(subscriber); err != nil {
		return nil, fmt.Errorf("操作失败")
	}
	return user, nil
}

func CancelSubscribeUser(userID uint, subscriberUserID uint) (*models.User, error) {
	var subscriber, user *models.User
	sql.DB.First(&subscriber, subscriberUserID)
	sql.DB.First(&user, userID)
	if subscriber == nil {
		return nil, fmt.Errorf("未找到关注人")
	}
	if user == nil {
		return nil, fmt.Errorf("未找到用户")
	}
	if err := sql.DB.Model(&user).Association("Subscribers").Delete(subscriber); err != nil {
		return nil, fmt.Errorf("关注不存在")
	}
	return user, nil
}

func GetUserSubscribersByID(userID uint) []models.User {
	var user *models.User
	sql.DB.Preload("Subscribers").Find(&user, userID)
	return user.Subscribers
}

func GetUserSubscribersCountByID(userID uint) int64 {
	var count int64
	sql.DB.Raw("select count(subscriber_id) from subscribes where user_id = ?",
		userID).Scan(&count)
	return count
}

func GetUserFollowersByID(userID uint) []models.User {
	var followers []models.User
	sql.DB.Raw("select * from users where id in"+
		"(select user_id from subscribes left join `users`"+
		"on `users`.id = subscriber_id "+
		"where subscriber_id = ?)", userID).Scan(&followers)
	return followers
}

func GetUserFollowersCountByID(userID uint) int64 {
	var count int64
	sql.DB.Raw("select count(user_id) from subscribes where subscriber_id = ?",
		&userID).Scan(&count)
	return count
}

func IsUserFollow(userID, anotherUserID uint) bool {
	if userID == anotherUserID {
		return false
	}
	var user *models.User
	sql.DB.Raw("select * from users where id in"+
		" (select user_id from subscribes where user_id = ? and subscriber_id = ?)",
		userID, anotherUserID).Scan(&user)
	return user != nil
}
