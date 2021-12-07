package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	var browser *rod.Browser
	if runtime.GOOS == "darwin" {
		u := launcher.New().Bin("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome").Headless(false).MustLaunch()
		browser = rod.New().ControlURL(u)
	} else if runtime.GOOS == "windows" {
		u := launcher.New().Bin("C:\\Users\\A\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe").Headless(false).MustLaunch()
		browser = rod.New().ControlURL(u)
	} else {
		browser = rod.New()
	}
	page := browser.MustConnect().MustPage("http://app1.nmpa.gov.cn/data_nmpa/face3/dir.html")
	page.MustWaitLoad()
	page.MustSearch("国产医疗器械产品").MustClick()
	page.MustWaitLoad()
	file, err := os.ReadFile("id.txt")
	if err != nil {
		log.Fatalln(err)
	}
	open, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	for _, s := range strings.Split(string(file), "\n") {
		s := strings.TrimSpace(s)
		if s != "" {
			content := getContentById(page, s)
			open.Write([]byte(fmt.Sprintf("%s\t%s\n", s, content["适用范围/预期用途"])))
			open.Sync()
		}
		time.Sleep(time.Second * 30)
	}
}

func getContentById(page *rod.Page, id string) map[string]string {
	for _, element := range page.MustElements("input") {
		property, err := element.Property("name")
		if err == nil && property.String() == "COLUMN180" {
			element.MustFocus().MustClick().MustInput(id).MustPress(input.Enter)
			e := proto.NetworkResponseReceived{}
			wait := page.WaitEvent(&e)
			wait()
			break
		}
	}
	page.MustWaitLoad()
	page.MustElement("#content > table:nth-child(2) > tbody > tr:nth-child(1) > td > p > a").MustClick()
	e := proto.NetworkResponseReceived{}
	wait := page.WaitEvent(&e)
	wait()
	page.MustWaitLoad()
	m := make(map[string]string)
	for _, element := range page.MustElements("#content > div > div > table:nth-child(1) tr") {
		elements := element.MustElements("td")
		if len(elements) >= 2 {
			m[elements[0].MustText()] = elements[1].MustText()
		}
	}
	return m
}
