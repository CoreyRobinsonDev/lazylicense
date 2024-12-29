package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const DOMAIN = "https://choosealicense.com"
const MAX_BOX_WIDTH = 45
const MAX_TERM_WIDTH = 102



func main() {
	licenses := GetLicenses()
	names := make([]string, 0, len(licenses))

	for _, license := range licenses {
		name := license.Name
		if license.AbbrName != "" {
			name += " (" + license.AbbrName + ")"
		}
		names = append(names, name)
	} 

	InitInput()
	moveNum := 0
	selection := HighlightSelection(moveNum, names)
	container, containerHeight := Container(LicenseDetails(licenses[0]), strings.Join(selection, "\n"))
	if TermWidth() <= MAX_TERM_WIDTH {
		fmt.Println(strings.Join(selection, "\n"))	
		MoveCursor("up", len(selection))
	} else {
		fmt.Println(container)
		MoveCursor("up", containerHeight)
	}
	MoveCursor("left", 99)

	for {
		// halts here
		moveNum += CalcInput()
		// clear last container from terminal
		for range containerHeight {
			fmt.Println(strings.Repeat(" ", int(TermWidth())))
		}
		MoveCursor("up", containerHeight)
		if moveNum < 0 {
			moveNum = len(licenses) - 1
		}
		idx := moveNum % len(licenses)
		selection := HighlightSelection(idx, names)
		if TermWidth() <= MAX_TERM_WIDTH {
			fmt.Println(strings.Join(selection, "\n"))
			MoveCursor("up", len(selection))
		} else {
			container, containerHeight := Container(LicenseDetails(licenses[idx]), strings.Join(selection, "\n"))
			fmt.Println(container)
			MoveCursor("up", containerHeight)
		}
		MoveCursor("left", 99)
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

