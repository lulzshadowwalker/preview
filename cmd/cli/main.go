package main

import (
	"log/slog"
	"os"

	"flag"

	preview "github.com/lulzshadowwalker/preview/pkg"
)

func main() {
	if len(os.Args) > 1 {
		slog.Error("invalid number of arguments")
		flag.Usage()
	}

	url := flag.String("url", "https://www.pinterest.com/pin/961166745447040917/", "url to get preview")
	flag.Parse()

	preview, err := preview.FromURL(*url)
	if err != nil {
		slog.Error("failed to get preview", "err", err)
	}

	slog.Info("preview", "preview", preview)
}
