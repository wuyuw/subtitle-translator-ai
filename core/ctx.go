package core

import (
	"path"

	"github.com/spf13/viper"
	"github.com/yanyiwu/gojieba"
)

type SvcContext struct {
	Jieba *gojieba.Jieba
}

func NewSvcContext() *SvcContext {
	return &SvcContext{
		Jieba: Newjieba(),
	}
}

func Newjieba() *gojieba.Jieba {
	dictDir := viper.GetString("jiebaDictDir")
	jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
	hmmPath := path.Join(dictDir, "hmm_model.utf8")
	userPath := path.Join(dictDir, "user.dict.utf8")
	idfPath := path.Join(dictDir, "idf.utf8")
	stopPath := path.Join(dictDir, "stop_words.utf8")
	return gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)
}
