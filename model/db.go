package model

import (
	"dousheng/config"
	mgo "gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

const (
	ColUser       = "user"
	ColVideo      = "video"
	ColAssessment = "assessment"
)

var (
	userMaxId  int = 0
	videoMaxId int = 0
	assMaxId   int = 0
)

var (
	DBName = config.Conf.Mongo.Name
)

func getCollection(col string) (collection *mgo.Collection, cls func()) {
	s := mongoSession.Clone()
	c := s.DB(DBName).C(col)
	return c, s.Close
}

func changeData(col string, query, update interface{}) error {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)

	err := c.Update(query, update)
	return err
}

func deleteOne(col string, query interface{}) error {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)
	return c.Remove(query)
}

func insertData(col string, query interface{}) error {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)
	return c.Insert(query)
}

func countData(col string, query interface{}) (int, error) {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)
	q := c.Find(query)

	return q.Count()

}

func userGet(query, selector interface{}) (User, error) {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(ColUser)

	user := User{}
	q := c.Find(query)
	if selector != nil {
		q.Select(selector)
	}
	err := q.One(&user)
	return user, err
}

func userList(query, selector interface{}) ([]User, error) {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(ColUser)

	list := []User{}
	q := c.Find(query)
	if selector != nil {
		q.Select(selector)
	}
	err := q.All(&list)
	return list, err
}

func initMaxId() {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(ColUser)

	var user User
	err := c.Find(nil).Sort("-id").One(&user)
	if err == nil {
		userMaxId = user.UserId + 1
	} else {
		userMaxId = 0
	}
}
