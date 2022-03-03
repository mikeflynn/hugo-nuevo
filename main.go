package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const VERSION = "0.2.1"

func getMonth() string {
	_, month, _ := time.Now().Date()

	return fmt.Sprintf("%02d", month)
}

func getYear() string {
	year, _, _ := time.Now().Date()

	return fmt.Sprintf("%d", year)
}

func getDay() string {
	_, _, day := time.Now().Date()

	return fmt.Sprintf("%d", day)
}

func getDir() string {
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get current directory.")
		os.Exit(1)
	}

	return curDir + "/"
}

func formatSlug(slug string, title string) string {
	if slug == "" {
		slug = title
	}

	slug = strings.ToLower(slug)

	badChars, _ := regexp.Compile(`[^\w\s]`)
	slug = badChars.ReplaceAllString(slug, "")

	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func readStdIn() string {
	body := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		text := scanner.Text()

		if len(text) != 0 {
			body = append(body, text)
		} else {
			break
		}
	}

	return strings.Join(body, "\n")
}

func readFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func downloadFile(source string, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Create Blank Dest File
	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Buffer URL
	resp, err := client.Get(source)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Save Buffer to Dest
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(source string, dest string) error {
	// Open Source File
	original, err := os.Open(source)
	if err != nil {
		return err
	}

	defer original.Close()

	// Create Dest File
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	new, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer new.Close()

	// Copy
	_, err = io.Copy(new, original)
	if err != nil {
		return err
	}

	return nil
}

func appendFile(path string, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return err
	}

	return nil
}

func writeFile(path string, text string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return err
	}

	return nil
}

func findImages(text string) [][]string {
	r, _ := regexp.Compile(`(?i)!\[[\w\-\s\._].\]\(([0-9a-z\-:_\.\/]+)\)`)
	return r.FindAllStringSubmatch(text, -1)
}

func updateMarkdownImages(text string, relativeDest string) (string, error) {
	matches := findImages(text)
	//fmt.Println(fmt.Sprintf("Found images: %d", len(matches)))

	if len(matches) > 0 {
		i := 0
		for _, v := range matches {
			i = i + 1
			destPath := fmt.Sprintf("%s/%s%s", relativeDest, fmt.Sprintf("%03d", i), filepath.Ext(v[1]))

			switch {
			case strings.HasPrefix(v[1], "http"): // Remote path
				if err := downloadFile(v[1], getDir()+destPath); err != nil {
					return "", err
				}

				text = strings.ReplaceAll(text, v[1], strings.Replace(destPath, "assets/", "", 1))
			case strings.HasPrefix(v[1], "/"): // Abs path
				if err := copyFile(v[0], getDir()+destPath); err != nil {
					return "", err
				}

				text = strings.ReplaceAll(text, v[1], strings.Replace(destPath, "assets/", "", 1))
			default: // Rel path
				continue
			}
		}
	} else {
		fmt.Println("No matches found for images.")
	}

	return text, nil
}

func findHeaderImage(text string) (string, string) {
	if len(text) > 750 {
		matches := findImages(text[0:750])
		if len(matches) > 0 {
			text = strings.ReplaceAll(text, matches[0][0], "")
			return text, matches[0][1]
		}
	}

	return text, ""
}

// "blog/#y/#m/#s.md" => "blog/%s/%s/%s.md" => blog/2022/02/test-post.md
func parsePathFormat(format string, title string, slug string) string {
	dataMap := map[string]func() string{
		"#m": getMonth,
		"#y": getYear,
		"#d": getDay,
	}

	for k, v := range dataMap {
		format = strings.ReplaceAll(format, k, v())
	}

	if strings.Contains(format, "#s") {
		format = strings.ReplaceAll(format, "#s", formatSlug(slug, title))
	}

	return format
}

func main() {
	// Flags
	input := flag.String("i", "", "The file path of the markdown file to post.")
	assetPath := flag.String("a", "assets/images", "The location of image files.")
	title := flag.String("t", "", "Set the title of the post.")
	slug := flag.String("s", "", "Set a custom slug; Defaults to the title.")
	editor := flag.String("e", "", "Path of editor command to open the resulting file with.")
	format := flag.String("p", "blog/#y/#m/#s.md", "The path of the new blog post.")
	publish := flag.Bool("publish", false, "If true, set the draft status to false.")
	version := flag.Bool("v", false, "Display version number.")
	flag.Parse()

	if *version == true {
		fmt.Println("hugo-nuevo: v" + VERSION)
		os.Exit(1)
	}

	if _, err := os.Stat(getDir() + "archetypes"); err != nil {
		fmt.Println("😕 This doesn't appear to be a Hugo directory. ")
		os.Exit(1)
	}

	body := ""
	if *input != "" {
		contents, err := readFile(*input)
		if err != nil {
			fmt.Println("Unable to read input markdown file. 🤬")
			os.Exit(1)
		}

		body = contents
	}

	body = strings.TrimSpace(body)

	/*
		if body == "" {
			fmt.Println("No post body found. 🤔")
			os.Exit(1)
		}
	*/
	if *title == "" {
		// Look for title in body
		r, _ := regexp.Compile(`(?im)\- # ([\w\s’'\.,;]+)$`)
		matches := r.FindAllStringSubmatch(body, 1)
		if len(matches) > 0 {
			foundTitle := strings.TrimSpace(matches[0][1])
			title = &foundTitle

			body = strings.Replace(body, matches[0][0]+"\n", "", 1)
		}
	}

	// Generate the stub file

	postPath := parsePathFormat(*format, *title, *slug)
	fullPostPath := "./content/" + postPath

	hugoNew := exec.Command("hugo", "new", postPath)
	if err := hugoNew.Run(); err != nil {
		fmt.Println(err.Error())
	}

	// Scan and update images
	fullAssetPath := fmt.Sprintf("%s/%s/%s", *assetPath, getYear(), formatSlug(*slug, *title))
	body, err := updateMarkdownImages(body, fullAssetPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Set header image
	body, headerImage := findHeaderImage(body)

	if headerImage != "" || *publish == true {
		data, err := readFile(fullPostPath)
		if err == nil && len(data) > 0 {
			// Set post to publish
			if *publish == true {
				data = strings.Replace(data, "draft: true", "draft: false", 1)
			}

			// Set header image
			if headerImage != "" {
				data = strings.Replace(data, "image: \"\"", fmt.Sprintf("image: \"%s\"", headerImage), 1)
			}

			writeFile(fullPostPath, strings.TrimSpace(data))
		}
	}

	// Save to file
	appendFile(fullPostPath, body)

	// Open it in an editor
	if *editor != "" {
		openEditor := exec.Command(*editor, fullPostPath)
		if err := openEditor.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}
}
