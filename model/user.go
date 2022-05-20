package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId    int           `bson:"id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Pwd       string        `bson:"password" json:"password"`
	FollCount int           `bson:"follow_count" json:"follow_count"`
	FansCount int           `bson:"follower_count" json:"follower_count"`
	Follower  []int         `bson:"follower" json:"follower"`
	Fans      []int         `bson:"fans" json:"fans"`
	FavVideo  []int         `bson:"fav_video" json:"fav_video"`
	IsFollow  bool          `bson:"is_follow" json:"is_follow"`
}

type UserInfo struct {
	UserId    int    `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	FollCount int    `bson:"follow_count" json:"follow_count"`
	FansCount int    `bson:"follower_count" json:"follower_count"`
	IsFollow  bool   `bson:"is_follow" json:"is_follow"`
}

//return user id if success or will be -1
func UserAdd(name, pwd string) (int, error) {
	flag, err := UserExist(name)
	if err != nil {
		return -1, err
	} else if flag {
		return -1, errors.New("user exist")
	}
	user := User{}
	user.Name = name
	user.Pwd = pwd
	user.UserId = userMaxId
	userMaxId++
	err = insertData(ColUser, user)
	return user.UserId, err
}

func UserExist(name string) (bool, error) {
	query := bson.M{
		"name": name,
	}
	num, err := countData(ColUser, query)
	if num != 0 {
		return true, err
	}
	return false, err
}
//database user info
func UserGetById(id int) (User, error) {
	query := bson.M{
		"id": id,
	}
	user, err := userGet(query, nil)
	return user, err
}
//user info back to app
func UserInfoById(id int) (UserInfo ,error){
	query := bson.M{
		"id": id,
	}
	user, err := userGet(query, nil)
	user_info:=UserInfo{};
	user_info.Name=user.Name;
	user_info.UserId=user.UserId;
	user_info.FansCount=user.FansCount;
	user_info.FollCount=user.FollCount;
	user_info.IsFollow=false;
	return user_info, err
}

func UserLogin(name, pwd string) (User, error) {
	query := bson.M{
		"name":     name,
		"password": pwd,
	}
	return userGet(query, nil)
}
