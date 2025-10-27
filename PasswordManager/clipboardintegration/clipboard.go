package clipboardintegration

import "fmt"

func WriteToClipBoard(text string) error {
	return clipboard.WriteAll(text)
}
