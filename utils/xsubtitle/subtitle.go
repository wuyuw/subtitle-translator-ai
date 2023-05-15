package xsubtitle

import (
	"fmt"

	"github.com/asticode/go-astisub"
)

type LineInfo struct {
	OringalText string
	Text        string
	VoiceName   string
	InlineStyle *astisub.StyleAttributes
	Style       *astisub.Style
}

type LineBatch struct {
	Max       int
	Lines     []LineBatchElem
	Paragraph string
}

type LineBatchElem struct {
	Index int
	Info  *LineInfo
}

// GetLineInfo 数据打平
func GetLineInfo(item *astisub.Item) *LineInfo {
	lineInfo := &LineInfo{}
	for _, line := range item.Lines {
		for _, lineItem := range line.Items {
			lineInfo.InlineStyle = lineItem.InlineStyle
			lineInfo.Style = lineItem.Style
			lineInfo.OringalText += fmt.Sprintf(" %s", lineItem.Text)
		}
		lineInfo.VoiceName = line.VoiceName
	}
	return lineInfo
}

func NewLine(lineInfo *LineInfo) []astisub.Line {
	return []astisub.Line{
		{
			Items: []astisub.LineItem{
				{
					InlineStyle: lineInfo.InlineStyle,
					Style:       lineInfo.Style,
					Text:        lineInfo.Text,
				},
			},
			VoiceName: lineInfo.VoiceName,
		},
	}
}

func (batch *LineBatch) IsFull() bool {
	return len(batch.Lines) >= batch.Max
}

func (batch *LineBatch) Append(index int, lineInfo *LineInfo) (err error) {
	if batch.IsFull() {
		return fmt.Errorf("batch is full")
	}
	batch.Lines = append(batch.Lines, LineBatchElem{index, lineInfo})
	// batch.Paragraph += fmt.Sprintf(", %s", lineInfo.OringalText)
	batch.Paragraph += fmt.Sprintf("\n%s", lineInfo.OringalText)
	return nil
}

func (batch *LineBatch) Reset() {
	batch.Lines = make([]LineBatchElem, 0)
	batch.Paragraph = ""
}
