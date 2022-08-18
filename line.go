package gocc

import (
	"fmt"
	"strings"
)

var maxBodyWidth int

func init() {
	maxBodyWidth = 0
}

type Line struct {
	Kind           LineKind
	Body           string
	Comment        string
	ExtendNL       int
	BlockId        string
	BlockSeparator string
	Nest           int
}

func NewSrcLine(body, comment string) *Line {
	if len(body) > maxBodyWidth {
		maxBodyWidth = len(body)
	}

	return &Line{
		Kind:    SourceCode,
		Body:    body,
		Comment: comment,
	}
}

func NewSeparator(blockId, separator string, nest int, nl bool) *Line {
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
	return &Line{
		Kind:    Comment,
		Comment: comment,
	}
}

func (l *Line) src() string {
	return fmt.Sprintf(
		"%s%s; %s%s",
		l.Body,
		strings.Repeat(" ", maxBodyWidth-len(l.Body)+20),
		l.Comment,
		strings.Repeat("\n", l.ExtendNL),
	)
}

func (l *Line) separator() string {
	return fmt.Sprintf("; %s %s%s", l.BlockId, strings.Repeat(l.BlockSeparator, l.Nest), strings.Repeat("\n", l.ExtendNL))
}

func (l *Line) comment() string {
	return fmt.Sprintf("; %s%s", l.Comment, strings.Repeat("\n", l.ExtendNL))
}

func (l *Line) String() string {
	switch l.Kind {
	case SourceCode:
		return l.src()
	case Comment:
		return l.comment()
	case Separator:
		return l.separator()
	}
	return ""
}
