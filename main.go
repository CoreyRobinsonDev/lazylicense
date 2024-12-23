package main

import (
	"fmt"
	"strings"
)

func main() {
	println(box("this is some text lol\nmore text"))
}

func box(text string) string {
	tl := "╭"
	tr := "╮"
	bl := "╰"
	br := "╯"
	h := "─"
	v := "│"
	pad := " "
	h_line := h
	text_seg := strings.Split(text, "\n")
	max_len := len(text_seg[0])
	paddings := make([]string, 0, 1)
	var out string

	for _, seg := range text_seg {
		if max_len < len(seg) {
			max_len = len(seg)
		}
	}

	for range max_len {
		h_line += h
	}

	for _, seg := range text_seg {
		r_pad := pad
		var diff uint
		if d := max_len - len(seg); d < 0 {
			diff = 0
		} else { diff = uint(d) }

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

	out += fmt.Sprintf("%s%s%s\n",
		bl, h_line, br,
	)

	return out
}


