package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const DOMAIN = "https://choosealicense.com"
const MAX_BOX_WIDTH = 45
const MAX_TERM_WIDTH = 102


func main() {
	width := TermWidth()
	licenses := GetLicenses()
	names := make([]string, 0, len(licenses))

	for _, license := range licenses {
		name := license.Name
		if license.AbbrName != "" {
			name += " (" + license.AbbrName + ")"
		}
		names = append(names, name)
	} 

	fmt.Printf("Move%s Quit%s\n", Purple("<jk|\u2B06\u2B07>"), Purple("<q>"))

	InitInput()
	moveNum := 0
	selection := HighlightOptions(moveNum, names)
	container, containerHeight := Container(LicenseDetails(licenses[0]), strings.Join(selection, "\n"))
	if width <= MAX_TERM_WIDTH {
		fmt.Println(strings.Join(selection, "\n"))	
		MoveCursor("up", len(selection))
	} else {
		fmt.Println(container)
		MoveCursor("up", containerHeight)
	}
	MoveCursor("left", 99)

	for {
		// halts here
		result := CalcInput()
		//
		if result > 1 {
			for range containerHeight {
				fmt.Println(strings.Repeat(" ", int(width)))
			}
			MoveCursor("up", containerHeight)
			idx := moveNum % len(licenses)
			fmt.Println(Blue(licenses[idx].Name))
			AddLicense(licenses[idx])
			break
		}
		moveNum += result
		// clear last container from terminal
		for range containerHeight {
			fmt.Println(strings.Repeat(" ", int(width)))
		}
		MoveCursor("up", containerHeight)
		if moveNum < 0 {
			moveNum = len(licenses) - 1
		}
		idx := moveNum % len(licenses)
		selection := HighlightOptions(idx, names)
		if width <= MAX_TERM_WIDTH {
			fmt.Println(strings.Join(selection, "\n"))
			MoveCursor("up", len(selection))
		} else {
			container, containerHeight := Container(LicenseDetails(licenses[idx]), strings.Join(selection, "\n"))
			fmt.Println(container)
			MoveCursor("up", containerHeight)
		}
		MoveCursor("left", 99)
	}
	// show cursor
	fmt.Fprint(os.Stdout, "\x1b[?25h")
}


func AddLicense(license License) {
	yearPlaceholders := []string{
		"<year>",
		"[year]",
		"[yyyy]",
	}
	namePlaceholders := []string{
		"<name of author>",
		"[name of copyright owner]",
		"[fullname]",
	}
	yearPlaceholder := ""
	namePlaceholder := ""

	for _, placeholder := range yearPlaceholders {
		if strings.Contains(license.Content, placeholder) {
			yearPlaceholder = placeholder
		}
	}


	for _, placeholder := range namePlaceholders {
		if strings.Contains(license.Content, placeholder) {
			namePlaceholder = placeholder
		}
	}

	if len(yearPlaceholder) != 0 {
		year := ""
		fmt.Printf("Enter year: \r")
		getYear:
		for {
			b := make([]byte, 1)
			os.Stdin.Read(b)
			switch b[0] {
			// 10 == \n
			case 10: 
				// clear year on invalid number
				if _, err := strconv.Atoi(year); err != nil {
					year = ""
					fmt.Printf("                                                          \r")
					fmt.Printf("Enter year: %s\r", year)
					continue
				}
				break getYear
			// 127 == backspace
			case 127:
				if len(year) == 0 {continue}
				fmt.Printf("                                                          \r")
				year = year[:len(year)-1]
			default:
				year += string(b)
			}
			fmt.Printf("Enter year: %s\r", year)
		}
		fmt.Println()
		license.Content = strings.ReplaceAll(license.Content, yearPlaceholder, year)
	}


	if len(namePlaceholder) != 0 {
		name := ""
		fmt.Printf("Enter name: \r")
		getName:
		for {
			b := make([]byte, 1)
			os.Stdin.Read(b)
			switch b[0] {
			// 10 == \n
			case 10: break getName
			// 127 == backspace
			case 127:
				if len(name) == 0 {continue}
				fmt.Printf("                                                          \r")
				name = name[:len(name)-1]
			default:
				name += string(b)
			}
			fmt.Printf("Enter name: %s\r", name)
		}
		fmt.Println()
		license.Content = strings.ReplaceAll(license.Content, namePlaceholder, name)
	}

	dat := []byte(license.Content)
	os.WriteFile("LICENSE", dat, 0644)
	dat, err := os.ReadFile("README.md")
	if err != nil {
		pwdCmd := exec.Command("pwd")
		dirBytes := Unwrap(pwdCmd.Output())
		fmt.Printf("No %s file was found in '%s'\n", Bold("README.md"),strings.Trim(string(dirBytes), " \t\n"))
		fmt.Println("Create README.md?")

		List([]string{"Yes", "No"}, func(selection any) {
			if selection == "Yes" {
				dir := strings.Split(string(dirBytes), "/")
				programName := dir[len(dir)-1]
				dat := []byte(fmt.Sprintf(
					"# %s\n\n## License\n[%s](./LICENSE)",
					programName,
					license.Name,
				))
				os.WriteFile("README.md", dat, 0644)
			}
		})
	} else {
		dat := []byte(fmt.Sprintf(
			"%s\n\n## License\n[%s](./LICENSE)",
			dat,
			license.Name,
		))
		os.WriteFile("README.md", dat, 0644)
	}
}

type License struct {
	Name string
	AbbrName string
	Description string
	Content string
	Permissions []string
	Conditions []string
	Limitations []string
}


func LicenseDetails(license License) string {
	detailsString := "Permissions\n\u2b24 " + strings.Join(license.Permissions, " \n\u2b24 ") + "\n\nConditions\n\u2b24 " + strings.Join(license.Conditions, " \n\u2b24 ") + "\n\nLimitations\n\u2b24 " + strings.Join(license.Limitations, " \n\u2b24 ")
	// TODO: calculate the width
	return Box(detailsString, MAX_BOX_WIDTH)
}


func GetLicenses() []License {
	res := Unwrap(http.Get(DOMAIN + "/licenses"))
	defer res.Body.Close()
	doc := Unwrap(goquery.NewDocumentFromReader(res.Body))
	licenses := make([]License, 0)

	doc.Find(".license-overview-name").Each(func(i int, s *goquery.Selection) {
		route, exists := s.Find("a").Attr("href")
		if !exists { handleErr("missing link on " + DOMAIN) }
		page := Unwrap(http.Get(DOMAIN + route))
		defer page.Body.Close()
		pageDoc := Unwrap(goquery.NewDocumentFromReader(page.Body))
		license := new(License)

		license.Name = strings.Trim(pageDoc.Find("h1").Text(), " \t\n")
		license.AbbrName = strings.Trim(pageDoc.Find(".license-nickname").Text(), " \t\n")
		var childOffset int
		if license.AbbrName == "" {
			childOffset = 1
		} else { childOffset = 2 }
		license.Description = strings.Trim(pageDoc.Find(fmt.Sprintf("div.license-body > p:nth-child(%d)", childOffset)).Text(), " \t\n")
		pageDoc.Find("ul.license-permissions > li").Each(func(i int, s *goquery.Selection) {
			license.Permissions = append(license.Permissions, strings.Trim(s.Text(), " \t\n"))
		})
		pageDoc.Find("ul.license-conditions > li").Each(func(i int, s *goquery.Selection) {
			license.Conditions = append(license.Conditions, strings.Trim(s.Text(), " \t\n"))
		})
		pageDoc.Find("ul.license-limitations > li").Each(func(i int, s *goquery.Selection) {
			license.Limitations = append(license.Limitations, strings.Trim(s.Text(), " \t\n"))
		})
		license.Content = pageDoc.Find("#license-text").Text()

		licenses = append(licenses, *license)
	})

	return licenses
}

