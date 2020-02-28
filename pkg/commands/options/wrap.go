package options

import "strings"

func wrap80(text string) string {
	return wrap(text, 80)
}

func wrap(text string, width int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	count := width - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > count {
			wrapped += "\n" + word
			count = width - len(word)
		} else {
			wrapped += " " + word
			count -= 1 + len(word)
		}
	}
	return wrapped
}
