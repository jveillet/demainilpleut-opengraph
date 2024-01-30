/*
Copyright © 2023 Jérémie Veillet <jeremie.veillet@gmail.com>
*/
package cmd

import (
	"fmt"
	"image/color"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	MARGIN_10 float64 = 10.0
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "demainilpleut's OpenGraph images generation",
	Long: `Opengraph is a CLI to generate opengraph images for blog posts.
it uses the command line arguments to write text on an image template.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := generate(title, author, file, labels, date, backgroundPath, logoPath)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
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
	// template background path
	backgroundPath string
	// Logo path
	logoPath string
)

func init() {
	rootCmd.AddCommand(generateCmd)

	// Local flags which will only run when this command is called directly
	generateCmd.Flags().StringVarP(&title, "title", "t", "", "post TITLE")
	generateCmd.Flags().StringVarP(&author, "author", "a", "", "post AUTHOR")
	generateCmd.Flags().StringVarP(&file, "file", "f", "", "output destination")
	generateCmd.Flags().StringVarP(&labels, "labels", "l", "", "post LABELS/TAGS")
	generateCmd.Flags().StringVarP(&date, "date", "d", "", "post DATE in YYYY-MM-DD format")
	generateCmd.Flags().StringVarP(&backgroundPath, "background_path", "b", "", "Background image temmplates path SRC")
	generateCmd.Flags().StringVarP(&logoPath, "logo_path", "o", "", "Logo image path SRC")

	// Required flags
	generateCmd.MarkFlagRequired("title")
	generateCmd.MarkFlagRequired("author")
	generateCmd.MarkFlagRequired("file")
	generateCmd.MarkFlagRequired("labels")
	generateCmd.MarkFlagRequired("date")
	generateCmd.MarkFlagRequired("background_path")
	generateCmd.MarkFlagRequired("logo_path")
}

// Opengraph image generation.
// This function generates an image based on a template, and output it in a destination file.
func generate(title string, author string, filePath string, tags string, date string, backgroundPath string, logoPath string) error {
	// Create new canvas of dimension 1200 x 630
	c := canvas.New(1200, 630)

	// Set the context
	ctx := setContext(c)

	// Load the background image
	backgroundImage, err := loadImage(backgroundPath)
	if err != nil {
		return err
	}

	// Draw the background image at coordinates
	ctx.DrawImage(0, 0, backgroundImage, canvas.DPMM(1.0))

	font, err := loadFonts()
	if err != nil {
		return err
	}

	// Base X and Y position on the canvas
	var posY float64 = 90
	var posX float64 = 50

	// Load the logo image
	// logoImage, err := loadImage(templatePath, "logo")
	logoImage, err := loadImage(logoPath)
	if err != nil {
		return err
	}

	// Draw the logo image at coordinates
	ctx.DrawImage(posX, posY, logoImage, canvas.DPMM(1.0))

	// Move the Y cursor to start the title
	posY = 150

	// Draw the multiline title text
	face := font.Face(200.0, canvas.Black, canvas.FontNormal)
	titleText := canvas.NewTextBox(face, title, 1100, 500, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(posX, posY, titleText)

	// Move the X and Y cursor for author details
	posY = 533
	posX = 280

	// Draw "by" in normal font
	face = font.Face(80.0, color.RGBA{91, 91, 102, 255}, canvas.FontNormal)
	byAuhtorLabel := canvas.NewTextBox(face, "by", 200, 50, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(posX, posY, byAuhtorLabel)

	// Move the X cursor to be next to the "by" label
	posX = posX + byAuhtorLabel.Bounds().W + MARGIN_10

	// Draw the author name at coordinates
	face = font.Face(80.0, color.RGBA{91, 91, 102, 255}, canvas.FontBold)
	authorLabel := canvas.NewTextBox(face, author, 200, 50, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(posX, posY, authorLabel)

	// Format the date to human readable format (US)
	// Example: January 2, 2006
	dateUS := formatDate(date)

	// Draw the "published on" text at coordinates
	posX = posX + authorLabel.Bounds().W + MARGIN_10
	face = font.Face(80.0, color.RGBA{91, 91, 102, 255}, canvas.FontNormal)
	publishedOnText := canvas.NewTextBox(face, fmt.Sprintf("published on %v", dateUS), 800, 50, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(posX, posY, publishedOnText)

	err = renderers.Write(filePath, c, canvas.DPI(24.0))
	if err != nil {
		logger.Error("Failed to create Opengraph image", "path", filePath)
		return err
	}

	logger.Info("Opengraph image successfully created", "path", filePath)

	return nil
}

// Creates a canvas context used to keep drawing state
func setContext(c *canvas.Canvas) *canvas.Context {
	ctx := canvas.NewContext(c)
	ctx.SetCoordSystem(canvas.CartesianIV)

	return ctx
}

// Load the system fonts
func loadFonts() (*canvas.FontFamily, error) {
	font := canvas.NewFontFamily("truetype")

	if err := font.LoadSystemFont("Arial", canvas.FontRegular); err != nil {
		return nil, err
	}

	if err := font.LoadSystemFont("Arial, Bold", canvas.FontBold); err != nil {
		return nil, err
	}

	return font, nil
}

// Load local images
func loadImage(path string) (*canvas.Image, error) {
	// Load the image data
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Create a new image object
	img, err := canvas.NewPNGImage(file)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

// Formats the date in a human readable way
// Example: January 2, 2006
func formatDate(date string) string {
	layoutUS := "January 2, 2006"
	layoutISO := "2006-01-02"
	t, _ := time.Parse(layoutISO, date)

	return t.Format(layoutUS)
}
