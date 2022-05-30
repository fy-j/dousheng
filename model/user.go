package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
)

//user info in database
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

//user info back to app
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

//judge user exist by name
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
func UserInfoById(id int) (UserInfo, error) {
	query := bson.M{
		"id": id,
	}
	user, err := userGet(query, nil)
	user_info := UserInfo{}
	user_info.Name = user.Name
	user_info.UserId = user.UserId
	user_info.FansCount = user.FansCount
	user_info.FollCount = user.FollCount
	user_info.IsFollow = false
	return user_info, err
}

//if user pwd or name wrong,it will basc error,or data be stored in user
func UserLogin(name, pwd string) (User, error) {
	query := bson.M{
		"name":     name,
		"password": pwd,
	}
	return userGet(query, nil)
}

//judge user whther is author fans
func UserIsFollowers(user_id, author_id int) (bool, error) {
	user, err := UserGetById(user_id)
	if err != nil {
		return false, err
	}
	for _, num := range user.Follower {
		if num == author_id {
			return true, nil
		}
	}
	return false, nil
}

//follow or cancel follow author,action 1 stand for follow,2 stand for cancel
func UserFollow(user_id, author_id, action int) error {
	if action != 1 && action != 2 {
		return errors.New("action must be 1 or 2")
	}
	if action == 1 {
		err := changeData(ColUser, bson.M{"id": user_id}, bson.M{"$addToSet": bson.M{"follower": author_id}})
		if err != nil {
			return err
		}
		err = changeData(ColUser, bson.M{"id": user_id}, bson.M{"$inc": bson.M{"follow_count": 1}})
		if err != nil {
			return err
		}
	} else {
		err := changeData(ColUser, bson.M{"id": author_id}, bson.M{"$pull": bson.M{"fans": user_id}})
		if err != nil {
			return err
		}
		err = changeData(ColUser, bson.M{"id": author_id}, bson.M{"$inc": bson.M{"follower_count": -1}})
		if err != nil {
			return err
		}
	}
	return nil
}

//user follow list
func UserFollowList(user_id int) ([]UserInfo, error) {
	user, err := userGet(bson.M{"id": user_id}, nil)
	if err != nil {
		return []UserInfo{}, err
	}
	var list []UserInfo
	list, err = userinfoList(bson.M{
		"id": bson.M{
			"$in": user.Follower,
		},
	}, nil)
	return list, err
}

//user follow list
func UserFansList(user_id int) ([]UserInfo, error) {
	user, err := userGet(bson.M{"id": user_id}, nil)
	if err != nil {
		return []UserInfo{}, err
	}
	var list []UserInfo
	list, err = userinfoList(bson.M{
		"id": bson.M{
			"$in": user.Fans,
		},
	}, nil)
	return list, err
}
