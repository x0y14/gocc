package gocc

import (
	"fmt"
	"golang.org/x/exp/utf8string"
	"strings"
	"unicode/utf8"
)

var maxBodyWidth int
var blockWidth int

func init() {
	maxBodyWidth = 0
	blockWidth = 20
}

type Line struct {
	Kind           LineKind
	Body           string
	Comment        string
	ExtendNL       int
	BlockId        string
	BlockSeparator string
	Nest           int
	Data           string
}

func NewSrcLine(body, comment string) *Line {
	width := utf8.RuneCountInString(body) + utf8.RuneCountInString("; ") + utf8.RuneCountInString(comment)
	if width > maxBodyWidth {
		maxBodyWidth = width
	}

	return &Line{
		Kind:    SourceCode,
		Body:    body,
		Comment: comment,
	}
}

func NewSeparator(blockId, separator string, nest int, nl bool) *Line {
	width := utf8.RuneCountInString("; ") + utf8.RuneCountInString(blockId) + (utf8.RuneCountInString(separator) * nest)
	if width > maxBodyWidth {
		maxBodyWidth = width
	}

	n := 0
	if nl {
		n = 1
	}

	return &Line{
		Kind:           Separator,
		BlockId:        blockId,
		BlockSeparator: separator,
		Nest:           nest,
		ExtendNL:       n,
	}
}

func NewComment(comment string) *Line {
	if 2+utf8.RuneCountInString(comment) > maxBodyWidth {
		maxBodyWidth = 2 + utf8.RuneCountInString(comment)
	}
	return &Line{
		Kind:    Comment,
		Comment: comment,
	}
}

//func (l *Line) src() string {
//	return fmt.Sprintf(
//		"%s%s; %s",
//		l.Body,
//		strings.Repeat(" ", maxBodyWidth-utf8.RuneCountInString(l.Body)+20),
//		l.Comment,
//	)
//}
//
//func (l *Line) separator() string {
//	return fmt.Sprintf("; %s %s", l.BlockId, strings.Repeat(l.BlockSeparator, l.Nest))
//}
//
//func (l *Line) comment() string {
//	return fmt.Sprintf("; %s", l.Comment)
//}

func (l *Line) String() string {
	var str string
	switch l.Kind {
	case SourceCode:
		whiteWidth := (maxBodyWidth + blockWidth) - (utf8.RuneCountInString(l.Body) + utf8.RuneCountInString(l.Comment) + 2) // utf8.RuneCountInString("; ") == 2
		str = fmt.Sprintf("%s%s; %s", l.Body, strings.Repeat(" ", whiteWidth), l.Comment)
	case Comment:
		str = "; " + l.Comment
	case Separator:
		str = fmt.Sprintf("; %s %s", l.BlockId, strings.Repeat(l.BlockSeparator, l.Nest))
	}

	if l.Nest == 1 && l.Kind == Separator {
		str += strings.Repeat("=", (maxBodyWidth+blockWidth)-utf8.RuneCountInString(str)) + " +"
	} else if l.Kind == SourceCode {
		str += " |"
	} else {
		if utf8string.NewString(str).IsASCII() {
			str += strings.Repeat(" ", (maxBodyWidth+blockWidth)-utf8.RuneCountInString(str)) + " |"
		}
	}

	str += strings.Repeat("\n", l.ExtendNL)

	return str
}
