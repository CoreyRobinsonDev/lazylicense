package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func InitInput() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func MoveCursor(dir string, amount int) {
	dirChar := ""	

	switch strings.ToLower(dir) {
	case "up": dirChar = "A"
	case "down": dirChar = "B"
	case "left": dirChar = "D"
	case "right": dirChar = "C"
	}
	fmt.Printf("\x1b[%d%s", amount, dirChar)
}

// -1 == up, 1 == down
func DetectMove() int {
	b := make([]byte, 3)
	os.Stdin.Read(b)
	if string(b[0]) == "j" || b[2] == 66 {
		return 1
	} else if string(b[0]) == "k" || b[2] == 65 {
		return -1
	}  

	return 0
}

func HighlightSelection(position int, options []string) []string {
	newOptions := make([]string, len(options))
	if position < 0 || position >= len(newOptions) { return newOptions }
	
	for i, el := range options {
		newOptions[i] = el
	}
	newOptions[position] = Blue(newOptions[position])

	return newOptions
}

func Container(left, right string) string {
	if len(left) == 0 {
		return right
	} else if len(right) == 0 {
		return left
	}
	leftArr := strings.Split(left, "\n")
	rightArr := strings.Split(right, "\n")
	leftWidth := utf8.RuneCountInString(leftArr[0])
	out := make([]string, 0)

	if len(leftArr) < len(rightArr) {
		for i, el := range rightArr {
			if i >= len(leftArr) {
				pad := strings.Repeat(" ", leftWidth)
				out = append(out, pad + el)
			} else {
				out = append(out, leftArr[i] + el)
			} 
		}
	} else {
		for i, el := range leftArr {
			if i >= len(rightArr) {
				out = append(out, el)
			} else {
				out = append(out, el + rightArr[i])
			} 
		}
	}

	return strings.Join(out, "\n")
}

// Returns the final box, width and height
func Box(text string) string  {
	tl := "╭"
	tr := "╮"
	bl := "╰"
	br := "╯"
	h := "─"
	v := "│"
	pad := " "
	h_line := h
	text_seg := strings.Split(text, "\n")
	max_len := utf8.RuneCountInString(text_seg[0])
	paddings := make([]string, 0, len(text_seg))
	var out string

	for _, seg := range text_seg {
		if max_len < utf8.RuneCountInString(seg) {
			max_len = utf8.RuneCountInString(seg)
		}
	}

	for range max_len {
		h_line += h
	}

	for _, seg := range text_seg {
		r_pad := pad
		var diff int
		if d := max_len - utf8.RuneCountInString(seg); d < 0 {
			diff = 0
		} else { diff = d }

		for range diff {
			r_pad += pad
		}

		paddings = append(paddings, r_pad)
	}

	h_line += h

	out = fmt.Sprintf("%s%s%s\n",
		tl, h_line, tr,
	)

	for i, seg := range text_seg {
		out += fmt.Sprintf("%s%s%s%s%s\n",
			v, pad, seg, paddings[i], v,
		)
	}

	out += fmt.Sprintf("%s%s%s",
		bl, h_line, br,
	)

	
	return out 
}


