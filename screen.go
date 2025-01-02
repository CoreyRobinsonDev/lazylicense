package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

func TermWidth() uint {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno)
    }
    return uint(ws.Col)
}

func InitInput() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// hide cursor
	fmt.Fprintf(os.Stdout, "\x1b[?25l")
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
// greater than 1 is a select
func CalcInput() int {
	b := make([]byte, 3)
	os.Stdin.Read(b)
	if string(b[0]) == "j" || b[2] == 66 {
		return 1
	} else if string(b[0]) == "k" || b[2] == 65 {
		return -1
	} else if b[0] == 9 || b[0] == 10 || b[0] == 32 {
		return 2
	}
	return 0
}

func HighlightOptions(position int, options []string) []string {
	newOptions := make([]string, len(options))
	if position < 0 || position >= len(newOptions) { return newOptions }
	
	for i, el := range options {
		newOptions[i] = el
	}
	newOptions[position] = Blue(newOptions[position])

	return newOptions
}

func List(options []string, onSelect func(selection any)) {
	moveNum := 0
	highlightOptions := HighlightOptions(moveNum, options)
	fmt.Println(strings.Join(highlightOptions, "\n"))
	MoveCursor("up", len(options))
	MoveCursor("left", 99)
	for {
		result := CalcInput()
		if result > 1 {
			idx := moveNum % len(options)
			MoveCursor("up", idx-2)
			fmt.Print(strings.Repeat(" ", int(TermWidth())))
			MoveCursor("left", 999)
			fmt.Println(options[idx])
			onSelect(options[idx])
			break
		}
		moveNum += result
		if moveNum < 0 {
			moveNum = len(options) - 1
		}
		idx := moveNum % len(options)
		selection := HighlightOptions(idx, options)
		fmt.Println(strings.Join(selection, "\n"))
		MoveCursor("up", len(selection))
		MoveCursor("left", 99)
	}
}

// Returns the container and max height
func Container(left, right string) (string, int) {
	leftHeight := strings.Count(left, "\n")
	rightHeight := strings.Count(right, "\n")
	if len(left) == 0 {
		return right, rightHeight
	} else if len(right) == 0 {
		return left, leftHeight
	}
	leftArr := strings.Split(left, "\n")
	rightArr := strings.Split(right, "\n")
	leftWidth := utf8.RuneCountInString(leftArr[0])
	out := make([]string, 0)

	if len(leftArr) < len(rightArr) {
		for i, el := range rightArr {
			if i >= len(leftArr) {
				pad := strings.Repeat(" ", leftWidth)
				out = append(out, pad + " " + el)
			} else {
				out = append(out, leftArr[i] + " " + el)
			} 
		}
	} else {
		for i, el := range leftArr {
			if i >= len(rightArr) {
				out = append(out, el)
			} else {
				out = append(out, el + " " + rightArr[i])
			} 
		}
	}

	return strings.Join(out, "\n"), max(leftHeight, rightHeight) + 1
}

func Box(text string, width ...int) string  {
	tl := "╭"
	tr := "╮"
	bl := "╰"
	br := "╯"
	h := "─"
	v := "│"
	pad := " "
	hLine := h
	textSeg := strings.Split(text, "\n")
	maxWidth := utf8.RuneCountInString(textSeg[0])
	if len(width) >= 1 {
		maxWidth = width[0]
	}
	paddings := make([]string, 0, len(textSeg))
	var out string

	for _, seg := range textSeg {
		if maxWidth < utf8.RuneCountInString(seg) {
			maxWidth = utf8.RuneCountInString(seg)
		}
	}

	for range maxWidth {
		hLine += h
	}

	for _, seg := range textSeg {
		r_pad := pad
		var diff int
		if d := maxWidth - utf8.RuneCountInString(seg); d < 0 {
			diff = 0
		} else { diff = d }

		for range diff {
			r_pad += pad
		}

		paddings = append(paddings, r_pad)
	}

	hLine += h

	out = fmt.Sprintf("%s%s%s\n",
		tl, hLine, tr,
	)

	for i, seg := range textSeg {
		out += fmt.Sprintf("%s%s%s%s%s\n",
			v, pad, seg, paddings[i], v,
		)
	}

	out += fmt.Sprintf("%s%s%s",
		bl, hLine, br,
	)

	return out 
}


