# `hugo-neuevo` 

A wrapper for the `hugo new post` that automates the process of importing posts from my chosen Markdown editor.

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

```