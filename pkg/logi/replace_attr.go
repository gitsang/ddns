package logi

import "log/slog"

type ReplaceAttrChain struct {
	ReplaceAttrs []ReplaceAttrFunc
}

func NewReplaceAttrChain() *ReplaceAttrChain {
	return &ReplaceAttrChain{
		ReplaceAttrs: []ReplaceAttrFunc{},
	}
}

func (r *ReplaceAttrChain) ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	for _, f := range r.ReplaceAttrs {
		a = f(groups, a)
	}
	return a
}

func (r *ReplaceAttrChain) Append(f ...ReplaceAttrFunc) {
	r.ReplaceAttrs = append(r.ReplaceAttrs, f...)
}
