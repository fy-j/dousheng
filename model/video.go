package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Video struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	AuthorID int           `bson:"author_id" json:"author_id"`
	VideoID  int           `bson:"id" json:"id"`
	PlayUrl  string        `bson:"play_url" json:"play_url"`
	CoverUrl string        `bson:"cover_url" json:"cover_url"`
	FavCount int           `bson:"favorite_count" json:"favorite_count"`
	ComCount int           `bson:"comment_count" json:"comment_count"`
	IsFav    bool          `bson:"is_favorite" json:"is_favorite"`
	Time     int64         `bson:"post_time" json:"post_time"`
	Title    string        `bson:"title" json:"title"`
}

type VideoInfo struct {
	VideoID  int      `bson:"id" json:"id"`
	PlayUrl  string   `bson:"play_url" json:"play_url"`
	CoverUrl string   `bson:"cover_url" json:"cover_url"`
	FavCount int      `bson:"favorite_count" json:"favorite_count"`
	ComCount int      `bson:"comment_count" json:"comment_count"`
	IsFav    bool     `bson:"is_favorite" json:"is_favorite"`
	Time     int64    `bson:"post_time" json:"post_time"`
	Author   UserInfo `bson:"author" json:"author"`
	Title    string   `bson:"title" json:"title"`
}

//get user self post video list
func VideoListByUserID(user_id int, times int64, limit int) ([]VideoInfo, error) {
	list, err := videoList(bson.M{"author_id": user_id, "post_time": bson.M{"$gt": times}}, nil, bson.M{"post_time": -1}, limit)
	if err != nil {
		return list, err
	}
	user, err := userGet(bson.M{"id": user_id}, nil)
	hashMap := make(map[int]bool)
	for _, num := range user.FavVideo {
		hashMap[num] = true
	}
	if err != nil {
		return list, err
	}
	for i := 0; i < len(list); i++ {
		if _, ok := hashMap[list[i].VideoID]; ok {
			list[i].IsFav = true
		}
	}
	return list, err
}

//get all video by time
func VideoList(time int64, limit int) ([]VideoInfo, error) {
	list, err := videoList(bson.M{"post_time": bson.M{"$gt": time}}, nil, bson.M{"post_time": -1}, limit)
	if err != nil {
		return list, err
	}
	return list, err
}

//get video meg by video id
func VideoMegByID(video_id int) (Video, error) {
	return videoGet(bson.M{
		"id": video_id,
	})
}

//is_fav is true if user_id's fav has video_id
func VideoMegByUserID(user_id, video_id int) (Video, error) {
	video, err := VideoMegByID(video_id)
	if err != nil {
		return video, err
	}
	user, err := UserGetById(user_id)
	if err != nil {
		return video, err
	}
	for _, fav := range user.FavVideo {
		if fav == video_id {
			video.IsFav = true
			break
		}
	}
	return video, err
}

//add video,return video id,if add wrong ,return -1 and error
func VideoAdd(user_id int, coverUrl, playUrl, title string) (int, error) {
	var video Video
	video.AuthorID = user_id
	video.CoverUrl = coverUrl
	video.PlayUrl = playUrl
	video.Title = title
	video.Time = time.Now().Unix()
	video.VideoID = videoMaxId
	videoMaxId++
	err := insertData(ColVideo, video)
	if err != nil {
		return -1, err
	}
	return video.VideoID, err
}

//return whether the video is user fav video
func VideoIsFav(user_id, video_id int) (bool, error) {
	user, err := userGet(bson.M{"id": user_id}, nil)
	if err != nil {
		return false, err
	}
	for _, num := range user.FavVideo {
		if num == video_id {
			return true, nil
		}
	}
	return false, nil
}

//video fav action,action 1 stand for fav,2 stand for cancel
func VideoFavAction(user_id, video_id, action int) error {
	if action != 1 && action != 2 {
		return errors.New("action wrong")
	}
	if action == 1 {
		err := changeData(ColUser, bson.M{"id": user_id}, bson.M{"$addToSet": bson.M{"fav_video": video_id}})
		if err != nil {
			return err
		}
		err = changeData(ColVideo, bson.M{"id": video_id}, bson.M{"$inc": bson.M{"favorite_count": 1}})
	} else {
		err := changeData(ColUser, bson.M{"id": user_id}, bson.M{"$pull": bson.M{"fav_video": video_id}})
		if err != nil {
			return err
		}
		err = changeData(ColVideo, bson.M{"id": video_id}, bson.M{"$inc": bson.M{"favorite_count": -1}})
	}
	return nil
}

//get user fav video list
func VideoFavList(user_id int) ([]VideoInfo, error) {
	user, err := UserGetById(user_id)
	if err != nil {
		return []VideoInfo{}, err
	}
	list, err := videoList(
		bson.M{
			"id": bson.M{
				"$in": user.FavVideo,
			}},
		nil,
		bson.M{
			"post_time": -1,
		},
		1000)
	if err != nil {
		return list, err
	}
	for i := 0; i < len(list); i++ {
		list[i].IsFav = true
	}
	return list, err
}

//return whether the video author is followed by user
func VideoAuthorIsFollowed(user_id, video_id int) (bool, error) {
	user, err := userGet(bson.M{"id": user_id}, nil)
	video, err := VideoMegByID(video_id)
	if err != nil {
		return false, err
	}
	for _, num := range user.Follower {
		if num == video.AuthorID {
			return true, nil
		}
	}
	return false, nil
}
