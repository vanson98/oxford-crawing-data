package controllers

import (
	"data-crawing/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/playwright-community/playwright-go"
)

type OxfordCrawingController struct {
	Page            playwright.Page
	InputSearch     playwright.Locator
	SearchButton    playwright.Locator
	EngPhoneticWord playwright.Locator
}

func InitOxfordCrawingController(browser playwright.Browser) OxfordCrawingController {
	page, err := browser.NewPage()
	page.SetDefaultTimeout(3000)
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto("https://www.oxfordlearnersdictionaries.com")
	page.Pause()
	crawingController := OxfordCrawingController{
		Page:            page,
		InputSearch:     page.Locator("//input[@id='q']"),
		SearchButton:    page.Locator("//label[@id='search-btn']/input"),
		EngPhoneticWord: page.Locator("//div[@class='phons_n_am']/span"),
	}
	return crawingController
}

func (occ OxfordCrawingController) CrawWordPhonetic(w http.ResponseWriter, r *http.Request) {
	requestDecoder := json.NewDecoder(r.Body)
	requestDataModel := struct {
		RangeWord []string
	}{}
	requestDecoder.Decode(&requestDataModel)

	occ.InputSearch.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(10000),
	})
	occ.SearchButton.WaitFor()

	result := make([]models.WordPhoneticModel, 0)

	for _, v := range requestDataModel.RangeWord {
		occ.InputSearch.Fill(v)
		occ.SearchButton.Click()
		occ.EngPhoneticWord.First().WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(3000),
			State:   playwright.WaitForSelectorStateVisible,
		})

		phons, _ := occ.EngPhoneticWord.TextContent()

		result = append(result, models.WordPhoneticModel{
			Word:     v,
			Phonetic: phons,
		})

	}
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}
