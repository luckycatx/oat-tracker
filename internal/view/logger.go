package view

import "gioui.org/widget"

type logger struct {
	list    *widget.List
	content []string
}

func (l *logger) Write(p []byte) (n int, err error) {
	s := string(p)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	l.content = append(l.content, s)
	return len(p), nil
}
