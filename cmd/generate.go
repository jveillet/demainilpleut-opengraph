/*
Copyright © 2023 Jérémie Veillet <jeremie.veillet@gmail.com>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	dcanvas "github.com/jveillet/demainilpleut-opengraph/pkg"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "demainilpleut's OpenGraph images generation",
	Long: `Opengraph is a CLI to generate opengraph images for blog posts.
it uses the command line arguments to write text on an image template.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := generate(title, author, output, date, backgroundPath, logoPath)
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
	output string
	// Post date (YYYY-MM-DD format)
	date string
	// template background path
	backgroundPath string
	// Logo path
	logoPath string
	/*go:embed assets/fonts*/
	// fs embed.FS

	//go:embed assets/fonts/Arial.ttf
	fontData []byte
)

func init() {
	rootCmd.AddCommand(generateCmd)

	// Local flags which will only run when this command is called directly
	generateCmd.Flags().StringVarP(&title, "title", "t", "", "post TITLE")
	generateCmd.Flags().StringVarP(&author, "author", "a", "", "post AUTHOR")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "output FILE")
	generateCmd.Flags().StringVarP(&date, "date", "d", "", "post DATE in YYYY-MM-DD format")
	generateCmd.Flags().StringVarP(&backgroundPath, "background_path", "b", "", "Background image temmplates path SRC")
	generateCmd.Flags().StringVarP(&logoPath, "logo_path", "l", "", "Logo image path SRC")

	// Required flags
	generateCmd.MarkFlagRequired("title")
	generateCmd.MarkFlagRequired("author")
	generateCmd.MarkFlagRequired("output")
	generateCmd.MarkFlagRequired("date")
	generateCmd.MarkFlagRequired("background_path")
	generateCmd.MarkFlagRequired("logo_path")
}

// Opengraph image generation.
// This function generates an image based on a template, and output it in a destination file.
func generate(title string, author string, output string, date string, backgroundPath string, logoPath string) error {
	// Create the canvas
	c := dcanvas.New(0, 0, 1200, 630)

	backgroundImage, err := loadImage(backgroundPath)
	if err != nil {
		log.Println("failed to load background image")
		log.Println("details : ", err)
		return err
	}

	// Add the background image to the canvas
	c.AddImage(c.Image(), c.Image().Bounds(), backgroundImage, backgroundImage.Bounds())

	// Load the logo image
	logoImage, err := loadImage(logoPath)
	if err != nil {
		log.Println("failed to load logo image")
		log.Println("details : ", err)
		return err
	}

	// Add the logo image to the canvas
	c.AddImage(c.Image(), image.Rect(50, 90, 50+logoImage.Bounds().Dx(), 90+logoImage.Bounds().Dy()), logoImage, logoImage.Bounds())

	// Load the system font "Arial"
	font, err := loadFont()
	if err != nil {
		log.Println("failed to load font")
		log.Println("details : ", err)
		return err
	}

	// Font configuration
	dc := freetype.NewContext()
	dc.SetDPI(72)
	dc.SetFont(font)
	dc.SetFontSize(80.0)
	dc.SetClip(c.Image().Bounds())
	dc.SetSrc(image.Black)
	dc.SetDst(c.Image())

	c.WithFreeTypeContext(dc)
	c.WithFont(font)

	// Draw the multiline title text
	err = c.DrawMultilineString(50, 200, title)
	if err != nil {
		log.Println("failed to draw multiline title")
		log.Println("details : ", err)
		return err
	}

	// Format the date to human readable format (US)
	// Example: January 2, 2006
	dateUS := formatDate(date)

	// Draw the author line
	dc.SetFontSize(28.0)
	authorLine := fmt.Sprintf("by @%v published on %v", author, dateUS)
	authorW, _ := c.MeasureText(28, authorLine)
	totalWidth := c.Image().Bounds().Dx()
	authorX := (totalWidth - int(authorW)) / 2

	err = c.DrawString(authorX, c.Image().Bounds().Dy()-50, authorLine)
	if err != nil {
		log.Println("failed to draw author line")
		log.Println("details : ", err)
		return err
	}

	err = saveImage(output, c.Image())
	if err != nil {
		log.Println("failed to create opengraph image")
		log.Println("path : ", output)
		log.Println("details : ", err)
		return err
	}

	log.Println("opengraph image successfully created")
	log.Println("path : ", output)

	return nil
}

// Load an embeded font (Arial.ttf)
// It returns a truetype Font object.
func loadFont() (*truetype.Font, error) {
	font, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return font, nil
}

// Load local images
// This function loads an image from the local filesystem.
// It returns an image.Image object.
// The path is the full path to the image file.
// Example: "/path/to/image.png"
func loadImage(path string) (image.Image, error) {
	// Load the image data
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new image object
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Save an image to the local filesystem
// This function saves an image to the local filesystem.
// The target image is an image.Image object.
// The file path is the full path to the image file.
// Example: "/path/to/image.png"
func saveImage(filePath string, target image.Image) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	err = png.Encode(f, target)
	if err != nil {
		return err
	}

	return nil
}

// Formats the date in a human readable way
// Example: January 2, 2006
func formatDate(date string) string {
	layoutUS := "January 2, 2006"
	layoutISO := "2006-01-02"
	t, _ := time.Parse(layoutISO, date)

	return t.Format(layoutUS)
}
