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

### 使用说明
1. 下载二进制文件及配置文件模板(此为Mac M1构建，其他系统架构自行打包): https://github.com/wuyuw/subtitle-translator-ai/releases/download/v1.0.0/sta-1.0.0-mac-arm64.zip
2. 更新配置文件
3. 执行脚本
```
$ ./sta -f config.yaml --inpath ./test.srt --outpath ./result.srt

# 命令行中指定的参数优先级更高，会覆盖配置文件中默认值
$ ./sta -f config.yaml --inpath ./test.srt --outpath ./result.srt --subject "network security" --batchSize=20
```

### 开发说明

1. 检出代码
```
$ git clone git@github.com:wuyuw/subtitle-translator-ai.git
```
2. 进入项目根目录
```
$ cd subtitle-translator-ai
```
3. 安装依赖项
```
$ go mod tidy
```
4. 构建
```
# Mac/Linux
$ go build -o sta .
# windows
$ go build -o sta.exe .
```

