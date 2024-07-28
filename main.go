package main

import (
	"data-crawing/controllers"
	"log"
	"net/http"

	"github.com/playwright-community/playwright-go"
)

func main() {
	// init playwrite
	browser := inithChromiumBrowser()
	oxfordCrawingController := controllers.InitOxfordCrawingController(browser)

	http.HandleFunc("/get-words-phonetic", oxfordCrawingController.CrawWordPhonetic)
	log.Fatal(http.ListenAndServe(":4500", nil))
}

func inithChromiumBrowser() playwright.Browser {
	err := playwright.Install()
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	return browser
}
