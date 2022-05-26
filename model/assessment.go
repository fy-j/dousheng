package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Assessment struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	AssID    int           `bson:"id" json:"id"`
	AuthorID int           `bson:"author_id" json:"author_id"`
	VideoID  int           `bson:"video_id" json:"video_id"`
	Content  string        `bson:"content" json:"content"`
	Time     int64         `bson:"date" json:"date"`
}

//ass info for app
type AssessmentInfo struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	AssID    int           `bson:"id" json:"id"`
	AuthorID int           `bson:"author_id" json:"author_id"`
	VideoID  int           `bson:"video_id" json:"video_id"`
	Content  string        `bson:"content" json:"content"`
	Time     int64         `bson:"date" json:"date"`
	Date     string        `bson:"create_date" json:"create_date"`
	User     User          `bson:"user" json:"user"`
}

//get video ass list
func AssListByVideoID(video_id int) ([]AssessmentInfo, error) {
	return assList(bson.M{
		"video_id": video_id,
	}, nil, bson.M{
		"date": -1,
	})
}

//add assessment
func AssAdd(user_id, video_id int, content string) error {
	var ass Assessment
	ass.AuthorID = user_id
	ass.Content = content
	ass.Time = time.Now().Unix()
	ass.AssID = assMaxId
	ass.VideoID = video_id
	assMaxId++
	err := insertData(ColAssessment, ass)
	if err != nil {
		return err
	}
	err = changeData(ColVideo, bson.M{"id": video_id}, bson.M{"$inc": bson.M{"comment_count": 1}})
	return err
}

//delete assessment
func AssDel(video_id, ass_id int) error {
	err := changeData(ColVideo, bson.M{"id": video_id}, bson.M{"$inc": bson.M{"comment_count": -1}})
	if err != nil {
		return err
	}
	return deleteOne(ColAssessment, bson.M{"id": ass_id})
}
