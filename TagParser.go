package goioc

import (
	"fmt"
	"github.com/gosrv/util"
	"reflect"
	"strconv"
)

/**
tag解析，使用go默认的tag解析
*/
type ITagParser interface {
	Parse(tag reflect.StructTag) map[string]string
}

var ITagParserType = reflect.TypeOf((*ITagParser)(nil)).Elem()

var TagParserHelper = struct {
	GetTagParser func(beanContainer IBeanContainer) ITagParser
}{
	GetTagParser: func(beanContainer IBeanContainer) ITagParser {
		tagParserInss := beanContainer.GetBeanByType(ITagParserType)
		switch len(tagParserInss) {
		case 0:
			return nil
		case 1:
			return tagParserInss[0].(ITagParser)
		default:
			parserTypes := ""
			for _, tp := range tagParserInss {
				parserTypes += fmt.Sprintf("%v ", reflect.TypeOf(tp))
			}
			util.Panic("expect only 1 tag parser, but find %v:[%v]", len(tagParserInss), parserTypes)
		}
		return nil
	},
}

type DefaultTagParser struct {
}

func NewTagParser() ITagParser {
	return &DefaultTagParser{}
}

func (this *DefaultTagParser) Parse(tag reflect.StructTag) map[string]string {
	s := make(map[string]string)
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		name := string(tag[:i])
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			if len(name) > 0 {
				s[name] = ""
			}
			break
		}
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			if len(name) > 0 {
				s[name] = ""
			}
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			panic(err)
		}
		s[name] = value
	}
	return s
}
