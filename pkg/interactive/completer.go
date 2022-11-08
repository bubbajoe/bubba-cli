package interactive

import (
	"strings"

	"github.com/bubbajoe/bubba-cli/pkg/util"
	"github.com/c-bata/go-prompt"
)

var namespaces = map[string]string{
	"search": "",
	"vsm":    "Vector Space Model Searching (tf-idf)",
	"env":    "",
	"exit":   "close the interactive session",
}

var namespaceSuggestions = map[string]func(string) []prompt.Suggest{}

func RegisterSuggestion(namespace string, suggestions []prompt.Suggest) bool {
	_, ok := namespaceSuggestions[namespace]
	if !ok {
		namespaceSuggestions[namespace] = func(s string) []prompt.Suggest {
			return suggestions
		}
	}
	return !ok
}

func RegisterSuggestionFunc(namespace string, suggestFunc func(string) []prompt.Suggest) bool {
	_, ok := namespaceSuggestions[namespace]
	if !ok {
		namespaceSuggestions[namespace] = suggestFunc
	}
	return !ok
}

func completer(in prompt.Document) []prompt.Suggest {
	s := util.MaptoSlice(namespaces,
		func(k string, v string) prompt.Suggest {
			if v != "" {
				return prompt.Suggest{Text: k, Description: v}
			}
			return prompt.Suggest{Text: k}
		},
	)
	m := s
	firstWord := strings.Split(in.TextBeforeCursor(), " ")[0]
	// fmt.Printf("'%s'-'%s'\n", in.TextBeforeCursor(), firstWord)
	if _, ok := namespaceSuggestions[firstWord]; ok {
		return prompt.FilterHasPrefix(
			namespaceSuggestions[firstWord](in.TextBeforeCursor()),
			in.GetWordBeforeCursor(), true)
	}
	return prompt.FilterHasPrefix(m, in.GetWordBeforeCursor(), true)
}
