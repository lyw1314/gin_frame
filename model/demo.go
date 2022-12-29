package model

import (
	"time"
)

type BlogArticle struct {
	TagId      int64  `json:"tag_id"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	ModifiedOn int64  `json:"modified_on"`
}

func (b BlogArticle) TableName() string {
	return "blog_article"
}
func (b *BlogArticle) GetList(condition map[string]interface{}) (list []BlogArticle, err error) {
	err = DbHandle["center"].Where(condition).Find(&list).Error
	return
}

func (b *BlogArticle) Add() error {
	nowTime := time.Now().Unix()
	b.ModifiedOn = nowTime
	err := DbHandle["center"].Save(b).Error
	return err
}
