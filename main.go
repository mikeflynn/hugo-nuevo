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

func main() {
	// Flags
	title := flag.String("t", "", "Set the title of the post.")
	slug := flag.String("s", "", "Set a custom slug; Defaults to the title.")
	editor := flag.String("e", "", "Path of editor command to open the resulting file with.")

	flag.Parse()

	if *title == "" {
		fmt.Println("The post title is required.")
		os.Exit(1)
	}

	postPathFormat := "blog/%s/%s/%s.md"
	fullPostPathFormat := "./content/" + postPathFormat

	hugoNew := exec.Command("hugo", "new", fmt.Sprintf(postPathFormat, getYear(), getMonth(), formatSlug(*slug, *title)))
	if err := hugoNew.Run(); err != nil {
		fmt.Println(err.Error())
	}

	fullPostPath := fmt.Sprintf(fullPostPathFormat, getYear(), getMonth(), formatSlug(*slug, *title))
	if hasStdIn() {
		appendText(fullPostPath, readStdIn())
	}

	if *editor != "" {
		openEditor := exec.Command(*editor, fullPostPath)
		if err := openEditor.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}
}
