package openai

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"subtitle-translator-ai/utils/xlog"

	openai "github.com/sashabaranov/go-openai"
)

var ErrorCodeMessage = map[int]string{
	401: "[OpenAI] 提供错误的API密钥 | Incorrect API key provided",
	403: "[OpenAI] 服务器拒绝访问，请稍后再试 | Server refused to access, please try again later",
	502: "[OpenAI] 错误的网关 |  Bad Gateway",
	503: "[OpenAI] 服务器繁忙，请稍后再试 | Server is busy, please try again later",
	504: "[OpenAI] 网关超时 | Gateway Time-out",
	500: "[OpenAI] 服务器繁忙，请稍后再试 | Internal Server Error",
}

type OpenAI struct {
	client *openai.Client
}

func NewClient(oenaiKey, proxy string) *OpenAI {
	openaiConfig := openai.DefaultConfig(oenaiKey)
	if proxy != "" {
		proxyUrl, err := url.Parse(fmt.Sprintf("http://%s", proxy))
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		openaiConfig.HTTPClient = &http.Client{
			Transport: transport,
		}
	} else {
		log.Println(xlog.Warn("未配置代理!"))
	}

	return &OpenAI{
		client: openai.NewClientWithConfig(openaiConfig),
	}
}

func (ai *OpenAI) Dial() error {
	_, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err
	}
	return nil
}

func (ai *OpenAI) Translate(message, parentMessageId, systemPrompt string) (messageId, content string, err error) {
	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   1000,
			Temperature: 0,
			TopP:        1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role: openai.ChatMessageRoleUser,
					// Content: "Okay, I understand. Please give me the content to be translated.",
					Content: "Translate from English to Simplified Chinese. Only the translated text can be returned. Only translate the text between <db01> and </db01>.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("<db01>%s</db01> =>", message),
				},
			},
		},
	)
	if err != nil {
		log.Print(xlog.Warn(fmt.Printf("ChatCompletion error: %v\n", err)))
		return "", "", err
	}
	runeContent := []rune(resp.Choices[0].Message.Content)
	return resp.ID, string(runeContent[6 : len(runeContent)-7]), nil
}
