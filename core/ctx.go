package core

import "github.com/yanyiwu/gojieba"

type SvcContext struct {
	Jieba *gojieba.Jieba
}

func NewSvcContext() *SvcContext {
	return &SvcContext{
		Jieba: gojieba.NewJieba(),
	}
}
