package preview

import (
	"net/url"
	"github.com/PuerkitoBio/goquery"
)

type GenericScrapper struct {
}

func (gs GenericScrapper) isSatisfied(url url.URL) bool {
	return true
}

func (gs GenericScrapper) TitleScrappers() []Extractor {
	return []Extractor{
		{
			Selector: "meta[property=\"og:title\"]",
			Handler: func(s *goquery.Selection) (*string, error) {
				content, ok := s.Attr("content")
				if !ok {
					return nil, ErrNotFound
				}

				return &content, nil
			},
		},
		{
			Selector: "title",
			Handler: func(s *goquery.Selection) (*string, error) {
				text := s.Text()
				if text == "" {
					return nil, ErrNotFound
				}

				return &text, nil
			},
		},
		{
			Selector: "h1",
			Handler: func(s *goquery.Selection) (*string, error) {
				text := s.Text()
				if text == "" {
					return nil, ErrNotFound
				}

				return &text, nil
			},
		},
	}
}

func (gs GenericScrapper) DescriptionScrappers() []Extractor {
	return []Extractor{
		{
			Selector: "meta[property=\"og:description\"]",
			Handler: func(s *goquery.Selection) (*string, error) {
				content, ok := s.Attr("content")
				if !ok {
					return nil, ErrNotFound
				}

				return &content, nil
			},
		},
		{
			Selector: "meta[name=\"description\"]",
			Handler: func(s *goquery.Selection) (*string, error) {
				content, ok := s.Attr("content")
				if !ok {
					return nil, ErrNotFound
				}

				return &content, nil
			},
		},
		{
			Selector: "h1 + p",
			Handler: func(s *goquery.Selection) (*string, error) {
				text := s.Text()
				if text == "" {
					return nil, ErrNotFound
				}

				return &text, nil
			},
		},
	}
}

func (gs GenericScrapper) ImageScrappers() []Extractor {
	return []Extractor{
		{
			Selector: "meta[property=\"og:image\"]",
			Handler: func(s *goquery.Selection) (*string, error) {
				content, ok := s.Attr("content")
				if !ok {
					return nil, ErrNotFound
				}

				return &content, nil
			},
		},
		{
			Selector: "img",
			Handler: func(s *goquery.Selection) (*string, error) {
				content, ok := s.Attr("src")
				if !ok {
					return nil, ErrNotFound
				}

				return &content, nil
			},
		},
		{
			Selector: "video[poster]",
			Handler: func(s *goquery.Selection) (*string, error) {
				src, ok := s.Attr("poster")
				if !ok {
					return nil, ErrNotFound
				}

				return &src, nil
			},
		},
	}
}
