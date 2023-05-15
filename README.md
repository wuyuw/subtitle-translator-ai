## SUBTITLE-TRANSLATOR-AI

基于OpenAI的离线字幕文件翻译脚本
```
$ sta -h
Usage:
  sta [flags]

Flags:
  -b, --batchSize int      每次翻译的字幕行数 (default 10)
  -f, --config string      配置文件
  -d, --endpoints string   指定语义完整的字幕序号集合(10,25,40)
  -e, --engine string      翻译引擎 (default "Google")
  -h, --help               help for sta
  -i, --inpath string      输入字幕文件路径
  -l, --language string    目标语言 (default "Chinese")
  -k, --openaiKey string   openai API key
  -o, --outpath string     输出字幕文件路径
  -p, --proxy string       代理配置(127.0.0.1:51837)
  -s, --subject string     目标语言 (default "movie")
```