# demainilpleut-opengraph

demainilpleut.dev's Opengraph image generator.

It is a cli tool written in Go that can be called from the main project to generate the Opengraph images for demainilpleut's articles.

This project was born out of the frustration of doing a similar thing with a node script, that was using a lot of dependecies, constantly
updated and breaking my build.

By rewriting this tool with Go, I can generate a binary and let it run for YEARS without touching it.

## Prerequisites

* go >= 1.21

The `Arial` and `Arial, Bold` system fonts are used by default, so make sure they are available on your system.

## Installation

Clone the repository.

```bash
git clone https://github.com/jveillet/demainilpleut-opengraph.git
```

Build the binary.

```bash
go build -o opengraph
```

Building without cgo (disables calling C code (import "C"))

```bash
CGO_ENABLED=0 go build -o opengraph
```

Or run via source.

```sh
go run main.go
```

## Usage

```bash
$ opengraph -h
A CLI tool to manipulate opengraph images for demainilpleut.dev

Usage:
  opengraph [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    demainilpleut's OpenGraph images generation
  help        Help about any command
  version     print Opengraph version

Flags:
  -h, --help     help for opengraph
  -t, --toggle   Help message for toggle

Use "opengraph [command] --help" for more information about a command.
```

### Generate

```bash
$ opengraph generate -h
Opengraph is a CLI to generate opengraph images for blog posts.
it uses the command line arguments to write text on an image template.

Usage:
  opengraph generate [flags]

Flags:
  -a, --author string            post AUTHOR
  -b, --background_path string   Background image temmplates path SRC
  -d, --date string              post DATE in YYYY-MM-DD format
  -f, --file string              output destination
  -h, --help                     help for generate
  -l, --labels string            post LABELS/TAGS
  -o, --logo_path string         Logo image path SRC
  -t, --title string             post TITLE
```

**Example:**

Via source

```bash
go run main.go generate -a johndoe -d 1970-01-01 -f ./dist/out.png -l "tag1, tag2, tag3" -t "The quick brown fox jumps over the lazy dog" -b "dist/background.png" -o "dist/logo.png"
```

Via binary

```bash
./opengraph generate -a johndoe -d 1970-01-01 -f ./dist/out.png -l "tag1, tag2, tag3" -t "The quick brown fox jumps over the lazy dog" -b "dist/background.png" -o "dist/logo.png"
```

## Breaking changes

Since version `2.0.0`:

Loading of images path from the environment is deprecated.

Image templates doesn't rely on environment variables anymore.

You need to pass the full path to the images via the arguments to the `generate` command:

* `-b` for the path to the background image.
* `-o` for the path to the logo.

As such, images have been removed from the repo and depend now on images external to this tool.

Using font files from the `config/` directory is deprecated, this tool uses the available System Fonts now.

## Can I use this ?

Well, it is specifically tailored to run for demainilpleut.dev, so it might not be usefull to you.

But if you want to explore code, or adapt it to your own needs, knock yourself out 😄️.
