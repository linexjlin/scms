package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func markdown2html(markdown string, title string) (string, error) {
	// Define the pandoc command with arguments
	pandocCmd := exec.Command(`./bin/pandoc.exe`, "--metadata", fmt.Sprintf("title=%s", title), "-B", "./template/h.html", "-A", "./template/f.html", "--katex")

	// Set the markdown content as the input to pandoc
	pandocCmd.Stdin = bytes.NewBufferString(markdown)

	// Capture the output of the pandoc command
	var out bytes.Buffer
	pandocCmd.Stdout = &out

	// Run the pandoc command
	if err := pandocCmd.Run(); err != nil {
		return "", fmt.Errorf("error running pandoc: %w", err)
	}

	return out.String(), nil
}
