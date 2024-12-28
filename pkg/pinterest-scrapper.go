package preview

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PinterestScrapper struct {
	gs GenericScrapper
}

func (ps PinterestScrapper) isSatisfied(url url.URL) bool {
	return strings.Contains(url.Host, "pinterest")
}

func (ps PinterestScrapper) TitleScrappers() []Extractor {
  extractors := []Extractor{}

  return append(extractors, ps.gs.TitleScrappers()...)
}

func (ps PinterestScrapper) DescriptionScrappers() []Extractor {
  extractors := []Extractor{
		{
			Selector: ".X8m.zDA.IZT.tBJ.dyH.iFc.j1A.swG",
			Handler: func(s *goquery.Selection) (*string, error) {
				text := s.Text()
				if text == "" {
					return nil, ErrNotFound
				}

				return &text, nil
			},
		},
	}

  return append(extractors, ps.gs.DescriptionScrappers()...)
}

func (ps PinterestScrapper) ImageScrappers() []Extractor {
  extractors := []Extractor{
		{
			Selector: ".hCL.kVc.L4E.MIw.N7A.XiG",
			Handler: func(s *goquery.Selection) (*string, error) {
				content, ok := s.Attr("src")
				if !ok {
					return nil, ErrNotFound
				}

				return &content, nil
			},
		},
	}

  return append(extractors, ps.gs.ImageScrappers()...)
}
