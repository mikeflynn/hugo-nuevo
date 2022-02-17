# `hugo-neuevo` 

A wrapper for the `hugo new post` that automates the process of importing posts from my chosen Markdown editor.

## Goals

The point of this project was to be able to export a markdown file from my chosen editor and quickly post it to my Hugo blog with as much automation as possible. The point of this project is to solve my needs, but I've tried to make it as generic as possible so others may benefit. 

## Installation

To install this on your computer, you can use the `go install` command:

```bash
> go install github.com/mikeflynn/hugo-nuevo@latest
```

## Usage

```bash
❯ ./hugo-nuevo -h
Usage of ./hugo-nuevo:
  -a string
      The location of image files. (default "assets/images")
  -e string
      Path of editor command to open the resulting file with.
  -i string
      The file path of the markdown file to post.
  -p string
      The path of the new blog post. (default "blog/#y/#m/#s.md")
  -s string
      Set a custom slug; Defaults to the title.
  -t string
      Set the title of the post.
  -v  Display version number.
```

## Example

```bash
> ./hugo-nuevo -i "~/Desktop/import_file.md" -a "assets/images/post"
```