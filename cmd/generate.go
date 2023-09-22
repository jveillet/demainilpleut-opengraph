/*
Copyright © 2023 Jérémie Veillet <jeremie.veillet@gmail.com>
*/
package cmd

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
	"github.com/spf13/cobra"
)

const (
	LINE_HEIGHT       float64 = 1.5
	TEXT_MARGIN_RIGHT float64 = 50.0
	TEXT_MARGIN_TOP           = 120.0
	AUTHOR_POS_X              = 90
	AUTHOR_POS_Y              = 560
	PUBLISHED_ON_TEXT         = "published on"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "demainilpleut's OpenGraph images generation",
	Long: `Opengraph is a CLI to generate opengraph images for blog posts.
it uses the command line arguments to write text on an image template.`,
	Run: func(cmd *cobra.Command, args []string) {
		generate(title, author, file, labels, date)
	},
}

var (
	// Post title
	title string
	// Post author
	author string
	// Destination file (full path)
	file string
	// Labels / tags (list of comma-separated string)
	labels string
	// Post date (YYYY-MM-DD format)
	date string
)

func init() {
	rootCmd.AddCommand(generateCmd)

	// Local flags which will only run when this command is called directly
	generateCmd.Flags().StringVarP(&title, "title", "t", "", "post TITLE")
	generateCmd.Flags().StringVarP(&author, "author", "a", "", "post AUTHOR")
	generateCmd.Flags().StringVarP(&file, "file", "f", "", "output destination")
	generateCmd.Flags().StringVarP(&labels, "labels", "l", "", "post LABELS/TAGS")
	generateCmd.Flags().StringVarP(&date, "date", "d", "", "post DATE in YYYY-MM-DD format")

	// Required flags
	generateCmd.MarkFlagRequired("title")
	generateCmd.MarkFlagRequired("author")
	generateCmd.MarkFlagRequired("file")
	generateCmd.MarkFlagRequired("labels")
	generateCmd.MarkFlagRequired("date")
}

// Opengraph image generation.
// This function generates an image based on a template, and output it in a destination file.
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

	titleHeight := calculateStringHeight(dc, title, maxWidth)

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
		log.Fatalf("[ERROR] The OpenGraph image %s was not created.", filePath)
		return err
	}

	log.Printf("The OpenGraph image %s was created.", filePath)

	return nil
}

// Calculate a multiline string height.
// A maxWidth argument is provided to wrap the string if it goes over a value.
func calculateStringHeight(dc *gg.Context, text string, maxWidth float64) float64 {
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
