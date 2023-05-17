package core

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"subtitle-translator-ai/translator"
	"subtitle-translator-ai/utils/xlog"
	"subtitle-translator-ai/utils/xsubtitle"

	"github.com/asticode/go-astisub"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/viper"
)

const SYSTEM_PROMPT = `
You are a professional translation engine, 
You specialize in the field of %s. 
please translate the text into a colloquial, professional, elegant and fluent content, 
without the style of machine translation. You must only translate the text content, never interpret it.
Target language: %s
`

var SvcCtx *SvcContext

func Run() {
	language := viper.GetString("language")
	subject := viper.GetString("subject")
	engine := viper.GetString("engine")
	jiebaDictDir := viper.GetString("jiebaDictDir")
	proxy := viper.GetString("proxy")
	inPath := viper.GetString("inpath")
	outPath := viper.GetString("outpath")

	batchSize := viper.GetInt("batchSize")
	endpoints := viper.GetString("endpoints")

	log.Println(xlog.Info("配置读取成功！"))
	log.Println(xlog.Info(fmt.Sprintf("翻译引擎: %s", engine)))
	log.Println(xlog.Info(fmt.Sprintf("代理地址: %s", proxy)))
	log.Println(xlog.Info(fmt.Sprintf("结巴分词字典: %s", jiebaDictDir)))
	log.Println(xlog.Info(fmt.Sprintf("目标语言: %s", language)))
	log.Println(xlog.Info(fmt.Sprintf("主题场景: %s", subject)))
	log.Println(xlog.Info(fmt.Sprintf("批次行数: %d", batchSize)))
	log.Println(xlog.Info(fmt.Sprintf("标记结束行: %s", endpoints)))
	log.Println(xlog.Info(fmt.Sprintf("原字幕文件: %s", inPath)))
	log.Println(xlog.Info(fmt.Sprintf("输出字幕文件: %s", outPath)))

	endpointLines := make(map[int]struct{}, 0)
	if endpoints != "" {
		for _, i := range strings.Split(endpoints, ",") {
			lineNo, err := strconv.Atoi(i)
			if err != nil {
				continue
			}
			endpointLines[lineNo] = struct{}{}
		}
	}

	// 初始化上下文
	SvcCtx = NewSvcContext()
	defer func() {
		SvcCtx.Jieba.Free()
	}()

	// 输入输出文件校验
	if !DoesFileExist(inPath) {
		log.Fatal(xlog.Fatal(fmt.Sprintf("原字幕文件不存在: %s", inPath)))
	}
	inFileInfo, _ := os.Stat(inPath)
	if inFileInfo.Size() == 0 {
		log.Fatal(xlog.Fatal("原字幕文件为空"))
	}

	if DoesFileExist(outPath) {
		log.Fatal(xlog.Fatal(fmt.Sprintf("输出字幕文件已存在: %s", outPath)))
	}

	if path.Ext(inPath) != ".srt" || path.Ext(outPath) != ".srt" {
		log.Fatal(xlog.Fatal("当前仅支持.srt格式的字幕文件"))
	}

	var transEngine *translator.Translator
	if engine == "OpenAI" {
		transEngine = translator.MustNewOpenAITranslator(viper.GetString("openaiKey"), viper.GetString("proxy"))
	} else if engine == "Google" {
		log.Println(xlog.Warn("Google翻译暂未上线"))
		os.Exit(1)
	} else {
		log.Fatalf("不支持的翻译引擎: %s", xlog.Fatal(engine))
	}

	s, err := astisub.OpenFile(inPath)
	if err != nil {
		log.Fatal(xlog.Fatal(err))
	}
	log.Println(xlog.Info(fmt.Sprintf("字幕总条数: %d", len(s.Items))))

	systemPrompt := fmt.Sprintf(SYSTEM_PROMPT, subject, language)
	// 进度条
	bar := pb.StartNew(len(s.Items))

	batch := &xsubtitle.LineBatch{Max: batchSize}
	parentMessageId := ""
	for i := 0; i < len(s.Items); i++ {
		lineInfo := xsubtitle.GetLineInfo(s.Items[i])
		if err := batch.Append(i, lineInfo); err != nil {
			log.Fatal(xlog.Fatal(err))
		}
		_, isEndpoint := endpointLines[i+1]
		if !batch.IsFull() && !isEndpoint {
			continue
		}
		// 翻译
		msgId, content, err := transEngine.Client.Translate(batch.Paragraph, parentMessageId, systemPrompt)
		if err != nil {
			log.Panicln(xlog.Warn(
				fmt.Sprintf("第%d~%d条字幕翻译失败: %s", batch.Lines[0].Index, batch.Lines[len(batch.Lines)-1].Index, err)))
			bar.Add(len(batch.Lines))
			batch.Reset()
			continue
		}
		parentMessageId = msgId
		// 当前批次处理
		batchHandle(s, batch, content, bar)
		batch.Reset()
	}
	// 最后一批处理
	if len(batch.Lines) > 0 {
		_, content, err := transEngine.Client.Translate(batch.Paragraph, parentMessageId, systemPrompt)
		if err != nil {
			log.Panicln(xlog.Warn(
				fmt.Sprintf("第%d~%d条字幕翻译失败: %s", batch.Lines[0].Index, batch.Lines[len(batch.Lines)-1].Index, err)))
			bar.Add(len(batch.Lines))
			batch.Reset()
		} else {
			batchHandle(s, batch, content, bar)
			batch.Reset()
		}
	}
	bar.Finish()
	log.Println(xlog.Info("翻译完成，准备输出到文件"))
	outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(xlog.Fatal(err))
	}
	defer outFile.Close()
	if err := s.WriteToSRT(outFile); err != nil {
		log.Fatal(xlog.Fatal(err))
	}

}

func batchHandle(s *astisub.Subtitles, batch *xsubtitle.LineBatch, content string, bar *pb.ProgressBar) {
	wordsArr := jiebaSplit(content)
	offset := 0
	for i, elem := range batch.Lines {
		// 没到最后一行已经被分完了
		if offset >= len(wordsArr) {
			continue
		}
		// 原始行长度占比
		lenPCT := math.Ceil(float64(len(elem.Info.OringalText))/float64(len(batch.Paragraph))*100) / 100
		wordsLen := int(lenPCT * float64(len(wordsArr)))
		// 长度已超出
		if i == len(batch.Lines)-1 || offset+wordsLen >= len(wordsArr) {
			elem.Info.Text = strings.Join(wordsArr[offset:], "")
		} else {
			elem.Info.Text = strings.Join(wordsArr[offset:offset+wordsLen], "")
		}
		offset = offset + wordsLen
		s.Items[elem.Index].Lines = []astisub.Line{
			{
				Items: []astisub.LineItem{
					{
						InlineStyle: elem.Info.InlineStyle,
						Style:       elem.Info.Style,
						Text:        elem.Info.Text,
					},
				},
				VoiceName: elem.Info.VoiceName,
			},
		}
		bar.Increment()
	}
}

func jiebaSplit(s string) []string {
	words := make([]string, 0)
	if len(s) == 0 {
		return words
	}
	for _, word := range SvcCtx.Jieba.Cut(s, true) {
		if word == "\n" {
			continue
		}
		words = append(words, word)
	}
	return words
}
