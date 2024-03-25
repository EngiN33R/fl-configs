package semantic

import "github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/inireader"

type String struct {
	*Value
}

type StringOption func(i *String)

func StrOpts(opts ...ValueOption) StringOption {
	return func(i *String) {
		for _, opt := range opts {
			opt(i.Value)
		}
	}
}

func NewString(section *inireader.Section, key string, opts ...StringOption) *String {
	s := &String{Value: NewValue(section, key)}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *String) Get() string {
	if s.optional && len(s.section.ParamMap[s.key]) == 0 {
		return ""
	}
	return s.section.ParamMap[s.key][0].First.AsString()
}

func (s *String) Set(value string) {
	if s.isComment() {
		s.Delete()
	}

	processed_value := inireader.UniParseStr(value)
	if len(s.section.ParamMap[s.key]) == 0 {
		s.section.AddParamToStart(s.key, (&inireader.Param{IsComment: s.isComment()}).AddValue(processed_value))
	}
	// implement SetValue in Section
	s.section.ParamMap[s.key][0].First = processed_value
	s.section.ParamMap[s.key][0].Values[0] = processed_value
}

func (s *String) Delete() {
	delete(s.section.ParamMap, s.key)
	for index, param := range s.section.Params {
		if param.Key == s.key {
			s.section.Params = append(s.section.Params[:index], s.section.Params[index+1:]...)
		}
	}
}
