package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

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

func formatSlug(slug string, title string) string {
	if slug == "" {
		slug = title
	}

	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func hasStdIn() bool {
	file := os.Stdin
	fi, err := file.Stat()
	if err != nil {
		return false
	}

	size := fi.Size()
	if size > 0 {
		return true
	}

	return false
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

func appendText(path string, text string) error {
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

func downloadFile(source string, dest string) error {
	return nil
}

func findMarkdownImages(text string) ([]string, error) {
	return []string{}, nil
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
	title := flag.String("t", "", "Set the title of the post.")
	slug := flag.String("s", "", "Set a custom slug; Defaults to the title.")
	editor := flag.String("e", "", "Path of editor command to open the resulting file with.")
	format := flag.String("p", "blog/#y/#m/#s.md", "The path of the new blog post.")
	flag.Parse()

	if *title == "" {
		fmt.Println("The post title is required.")
		os.Exit(1)
	}

	// Generate the stub file

	postPath := parsePathFormat(*format, *title, *slug)
	fullPostPath := "./content/" + postPath

	hugoNew := exec.Command("hugo", "new", postPath)
	if err := hugoNew.Run(); err != nil {
		fmt.Println(err.Error())
	}

	// Fill in the body
	if hasStdIn() {
		appendText(fullPostPath, readStdIn())
	}

	// Scan and update images

	// Open it in an editor
	if *editor != "" {
		openEditor := exec.Command(*editor, fullPostPath)
		if err := openEditor.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}
}
