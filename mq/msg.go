package mq

import (
	"bytes"
	"encoding/gob"
	"log"
)

type PublishMsg struct {
	UserId   int    ` json:"user_id"`
	FileName string `json:"file_name"`
	Title    string `json:"title"`
}

func StructToBytes(v PublishMsg) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(v)
	if err != nil {
		log.Fatal(err)
	}
	return buffer.Bytes()
}

//Video反序列化
func BytesToStruct(byte []byte) PublishMsg {
	var msg PublishMsg
	decoder := gob.NewDecoder(bytes.NewReader(byte))
	err := decoder.Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
