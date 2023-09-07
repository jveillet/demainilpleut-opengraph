# demainilpleut-opengraph

demainilpleut.dev's Opengraph image generator.

It is a cli tool written in Go that can be called from the main project to generate the Opengraph images for demainilpleut's articles.

This project was born out of the frustration of doing a similar thing with a node script, that was using a lot of dependecies, constantly
updated and breaking my build.

By moving rewriting this tool with Go, I can generate a binary and let it run for YEARS without touching it.

## Prerequisites

* go >= 1.20

Fonts are available in the `fonts/` directory (Arial by default).
Image templates are in the `templates/` directory.

## Installation

Clone the repository.

```bash
git clone https://github.com/jveillet/demainilpleut-opengraph.git
```

Build the binary.

```bash
go build -o dist/opengraph
```

Building without cgo (disables calling C code (import "C"))

```bash
CGO_ENABLED=0 go build -o dist/opengraph
```

Or run via source.

```sh
go run main.go
```

## Configuration

Create a .env file at the root of the project, then add the two needed variables.

```bash
touch .env && OG_IMG_PATH=templates >> .env && OG_FONTS_PATH=fonts >> .env
```

## Usage

```bash
NAME:
   Opengraph - demainilpleut's OpenGraph images generation

USAGE:
   Opengraph [global options] command [command options] [arguments...]

VERSION:
   1.0.2

COMMANDS:
   generate, g  Generate an OpenGraph image
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Generate

```bash
NAME:
   Opengraph generate - Generate an OpenGraph image

USAGE:
   Opengraph generate [command options] [arguments...]

OPTIONS:
   --title TITLE, -t TITLE     The post TITLE
   --author AUTHOR, -a AUTHOR  The post AUTHOR
   --file PATH, -f PATH        Save the generated image in PATH
   --labels LABELS, -l LABELS  The post LABELS
   --date DATE, -d DATE        The post DATE in YYYY-MM-DD format
   --help, -h                  show help
```

## Can I use this ?

Well, it is specifically tailored to run for demainilpleut.dev, so it might not be usefull to you.

But if you want to explore code, or adapt it to your own needs, knock yourself out 😄️.
