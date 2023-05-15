package translator

import (
	"log"
	"subtitle-translator-ai/translator/openai"
	"subtitle-translator-ai/utils/xlog"
)

type TranslatorClient interface {
	Dial() error
	Translate(message, parentMessageId, systemPrompt string) (messageId, content string, err error)
}

type Translator struct {
	Engine string
	Client TranslatorClient
}

func MustNewOpenAITranslator(openaiKey, proxy string) *Translator {
	translator, err := NewOpenAITranslator(openaiKey, proxy)
	if err != nil {
		log.Fatal(xlog.Fatal(err))
	}
	return translator
}

func NewOpenAITranslator(openaiKey, proxy string) (*Translator, error) {
	client := openai.NewClient(openaiKey, proxy)
	if err := client.Dial(); err != nil {
		return nil, err
	}
	return &Translator{
		Engine: "OpenAI",
		Client: client,
	}, nil
}
