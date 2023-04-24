package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const (
	LINE_HEIGHT       float64 = 1.5
	TEXT_MARGIN_RIGHT float64 = 50.0
	TEXT_MARGIN_TOP           = 120.0
	AUTHOR_POS_X              = 90
	AUTHOR_POS_Y              = 560
	PUBLISHED_ON_TEXT         = "published on"
)

func main() {
	var title string
	var author string
	var filePath string
	var tags string
	var date string

	env := os.Getenv("BRIDGETOWN_ENV")
	if env == "" {
		env = "development"
	}

	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := &cli.App{
		Name:    "Opengraph",
		Usage:   "demainilpleut's OpenGraph images generation",
		Version: "1.0.1",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate an OpenGraph image",
				Action: func(cCtx *cli.Context) error {
					generate(cCtx.String("title"), cCtx.String("author"), cCtx.String("file"), cCtx.String("labels"), cCtx.String("date"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "title",
						Aliases:     []string{"t"},
						Usage:       "The post `TITLE`",
						Required:    true,
						Destination: &title,
					},
					&cli.StringFlag{
						Name:        "author",
						Aliases:     []string{"a"},
						Usage:       "The post `AUTHOR`",
						Required:    true,
						Destination: &author,
					},
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Usage:       "Save the generated image in `PATH`",
						Required:    true,
						Destination: &filePath,
					},
					&cli.StringFlag{
						Name:        "labels",
						Aliases:     []string{"l"},
						Usage:       "The post `LABELS`",
						Required:    true,
						Destination: &tags,
					},
					&cli.StringFlag{
						Name:        "date",
						Aliases:     []string{"d"},
						Usage:       "The post `DATE` in YYYY-MM-DD format",
						Required:    true,
						Destination: &date,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func generate(title string, author string, filePath string, tags string, date string) error {
	// Create the canvas context
	dc := gg.NewContext(1200, 630)

	// Load template
	backgroundImage, err := gg.LoadImage(fmt.Sprintf("%s/opengraph_image.png", os.Getenv("OG_IMG_PATH")))
	if err != nil {
		return err
	}

	// Use the template as a background
	dc.DrawImage(backgroundImage, 0, 0)

	// Add title text
	fontPath := filepath.Join(os.Getenv("OG_FONTS_PATH"), "Arial.ttf")
	fontPathBold := filepath.Join(os.Getenv("OG_FONTS_PATH"), "Arial_Bold.ttf")

	if err := dc.LoadFontFace(fontPath, 60); err != nil {
		return err
	}

	dc.SetColor(color.Black)

	maxWidth := float64(dc.Width()) - TEXT_MARGIN_RIGHT - TEXT_MARGIN_TOP

	dc.DrawStringWrapped(title, TEXT_MARGIN_RIGHT+1, TEXT_MARGIN_TOP+1, 0, 0, maxWidth, LINE_HEIGHT, gg.AlignLeft)

	titleHeight := textHeight(dc, title, maxWidth)

	// Add the tag icon
	tagImage, err := gg.LoadImage(fmt.Sprintf("%s/icon-tag.png", os.Getenv("OG_IMG_PATH")))
	if err != nil {
		return err
	}

	tagPositionY := int(TEXT_MARGIN_TOP) + int(titleHeight) + 30

	// Use the template as a background
	dc.DrawImage(tagImage, int(TEXT_MARGIN_RIGHT), tagPositionY)

	dc.SetColor(color.RGBA{133, 133, 133, 255})
	if err := dc.LoadFontFace(fontPath, 40); err != nil {
		return err
	}

	dc.DrawString(tags, float64(tagImage.Bounds().Dx())+60, float64(tagPositionY)+30)

	// Add author name
	if err := dc.LoadFontFace(fontPathBold, 30); err != nil {
		return err
	}

	dc.SetColor(color.Black)
	dc.DrawString(author, AUTHOR_POS_X, AUTHOR_POS_Y)

	authorWidth, _ := dc.MeasureString(author)

	// Add "published on"
	if err := dc.LoadFontFace(fontPath, 30); err != nil {
		return err
	}

	publishedOnPositionX := AUTHOR_POS_X + authorWidth + 10

	dc.DrawString(PUBLISHED_ON_TEXT, publishedOnPositionX, AUTHOR_POS_Y)

	publishedOnWidth, _ := dc.MeasureString(PUBLISHED_ON_TEXT)

	datePositionX := publishedOnPositionX + publishedOnWidth + 10

	// Add date locale
	layoutUS := "January 2, 2006"
	layoutISO := "2006-01-02"
	t, _ := time.Parse(layoutISO, date)
	parsedDate := t.Format(layoutUS)

	dc.DrawString(parsedDate, datePositionX, AUTHOR_POS_Y)

	// Save image
	err = dc.SavePNG(filePath)
	if err != nil {
		log.Printf("[ERROR] The OpenGraph image %s was not created.", filePath)
		return err
	}

	log.Printf("The OpenGraph image %s was created.", filePath)

	return nil
}

func textHeight(dc *gg.Context, text string, maxWidth float64) float64 {
	lines := dc.WordWrap(text, maxWidth)
	mls := ""
	for i, sl := range lines {
		mls = mls + sl
		if i != len(lines)-1 {
			mls = mls + "\n"
		}
	}

	_, textHeight := dc.MeasureMultilineString(mls, LINE_HEIGHT)

	return textHeight
}
