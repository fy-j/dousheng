package model

import (
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
	VideoID  int           `bson:"id" json:"id"`
	PlayUrl  string        `bson:"play_url" json:"play_url"`
	CoverUrl string        `bson:"cover_url" json:"cover_url"`
	FavCount int           `bson:"favorite_count" json:"favorite_count"`
	ComCount int           `bson:"comment_count" json:"comment_count"`
	IsFav    bool          `bson:"is_favorite" json:"is_favorite"`
	Time     int64         `bson:"post_time" json:"post_time"`
	Author   UserInfo      `bson:"author" json:"author"`
	Title    string        `bson:"title" json:"title"`
}
//get user self post video list
func VideoListByUserID(user_id int,times int64)([]VideoInfo,error){
	list,err:=videoList(bson.M{"author_id":user_id,"post_time":bson.M{"$gt":times}},nil,bson.M{"post_time":-1},100);
	if err!=nil{
		return list,err;
	}
	user,err:=userGet(bson.M{"id":user_id},nil);
	hashMap:=make(map[int]bool);
	for _,num:=range user.FavVideo{
		hashMap[num]=true;
	}
	if err!=nil{
		return list,err;
	}
	for i:=0;i<len(list);i++{
		if _,ok:=hashMap[list[i].VideoID];ok{
			list[i].IsFav=true;
		}
	}
	return list,err;
}
//get all video by time
func VideoList(time int64)([]VideoInfo,error){
	list,err:=videoList(bson.M{"post_time":bson.M{"$gt":time}},nil,bson.M{"post_time":-1},100);
	if err!=nil{
		return list,err;
	}
	return list,err;
}
//get video meg by video id
func VideoMegByID(video_id int)(Video,error){
	return videoGet(bson.M{
		"id":video_id,
	});
}
//add video,return video id,if add wrong ,return -1 and error
func VideoAdd(user_id int,coverUrl,playUrl,title string)(int,error){
	var video Video;
	video.AuthorID=user_id;
	video.CoverUrl=coverUrl;
	video.PlayUrl=playUrl;
	video.Title=title;
	video.Time= time.Now().Unix();
	video.VideoID=videoMaxId;
	videoMaxId++;
	err:=insertData(ColVideo,video);
	if err!=nil{
		return -1,err;
	}
	return video.VideoID,err;
}
