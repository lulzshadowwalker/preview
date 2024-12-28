package preview

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"log/slog"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

var (
	ErrNotFound = errors.New("not found")
)

var scrappers = []Scrapper{
	GenericScrapper{},
}

var pw *playwright.Playwright
var browser playwright.Browser

func init() {
	err := playwright.Install()
	if err != nil {
		slog.Error("failed to install playwright", "err", err)
		panic(err)
	}

	pw, err = playwright.Run()
	if err != nil {
		slog.Error("failed to run playwright", "err", err)
		panic(err)
	}

	browser, err = pw.Chromium.Launch()
	if err != nil {
		slog.Error("failed to launch browser", "err", err)
		panic(err)
	}

	slog.Info("browser launched")
}

type Preview struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}

func (p Preview) String() string {
	nullify := func(s *string) string {
		if s == nil {
			return "null"
		}
		return *s
	}

	return fmt.Sprintf("Title: %s\nDescription: %s\nImage: %s\n", nullify(p.Title), nullify(p.Description), nullify(p.Image))
}

type Scrapper interface {
	TitleScrappers() []Extractor
	DescriptionScrappers() []Extractor
	ImageScrappers() []Extractor
}

type Extractor struct {
	Selector string
	Handler  func(*goquery.Selection) (*string, error)
}

func (e Extractor) Extract(doc *goquery.Document) (*string, error) {
	s := doc.Find(e.Selector).First()
	if s.Length() == 0 {
		return nil, ErrNotFound
	}

	return e.Handler(s)
}

func FromURL(url string) (Preview, error) {
	ctx := context.WithValue(context.Background(), "url", url)

	page, err := browser.NewPage()
	if err != nil {
		return Preview{}, fmt.Errorf("failed to create new page: %v", err)
	}
	defer page.Close()

	if _, err := page.Goto(url); err != nil {
		return Preview{}, fmt.Errorf("failed to goto URL: %v", err)
	}

	page.WaitForSelector("h1")
	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateLoad,
		// Timeout: ,
	}); err != nil {
		return Preview{}, fmt.Errorf("failed to wait for load state: %v", err)
	}

	content, err := page.Content()
	if err != nil {
		return Preview{}, fmt.Errorf("failed to get page content: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return Preview{}, fmt.Errorf("failed to parse document: %v", err)
	}

	preview := Preview{}
	// scrapper := GenericScrapper{}
	scrapper := PinterestScrapper{}
	if preview.Title, err = Scrape(ctx, doc, scrapper.TitleScrappers()...); err != nil {
		slog.Warn("Title not found", "err", err)
	}

	if preview.Description, err = Scrape(ctx, doc, scrapper.DescriptionScrappers()...); err != nil {
		slog.Warn("Description not found", "err", err)
	}

	if preview.Image, err = Scrape(ctx, doc, scrapper.ImageScrappers()...); err != nil {
		slog.Warn("Image not found", "err", err)
	}

	if preview.Title == nil && preview.Description == nil && preview.Image == nil {
		return Preview{}, ErrNotFound
	}

	return preview, nil
}

func Scrape(ctx context.Context, doc *goquery.Document, extractors ...Extractor) (*string, error) {
	for _, e := range extractors {
		content, err := e.Extract(doc)
		if err != nil {
			continue
		}

		return content, nil
	}

	return nil, ErrNotFound
}

func Close() {
	pw.Stop()
	slog.Info("playwright stopped")

	browser.Close()
	slog.Info("browser closed")
}
