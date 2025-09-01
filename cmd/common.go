/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gocolly/colly"
)

type parameters struct {
	CurVersion string
	NewVersion string
	TempDir    string
	InstallDir string
	DlFileName string
	Extension  string
}

func isCacheValid(cacheFile string) bool {
	fCache, err := os.Open(cacheFile)
	if err != nil {
		// log.Fatal(err)
		return false
	}
	defer fCache.Close()
	// Check if cache file is empty
	cacheInfo, err := fCache.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if cacheInfo.Size() == 0 {
		return false
	}
	// Check if cache is still valid
	return maxCacheTime != 0.0 && time.Since(cacheInfo.ModTime()).Hours() <= maxCacheTime
}

func unCache(URL string) {
	sum := sha1.Sum([]byte(URL))
	hash := hex.EncodeToString(sum[:])
	dir := path.Join(cacheDir, hash[:2])
	filename := path.Join(dir, hash)
	if isCacheValid(filename) {
		return
	}
	fmt.Println("Deleting cached file.")
	if err := os.Remove(filename); err != nil {
		log.Fatal(err)
	}
}

func getCurrentVersion() {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		curVersion = "Go is not installed."
	} else {
		ver := strings.Split(string(out), " ")
		curVersion = strings.TrimPrefix(ver[2], "go")
	}
}

func scrapeLatestVersion() {
	_, err := os.Stat(cacheDir)
	if err == nil {
		unCache("https://go.dev/dl/")
	}
	c := colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.AllowURLRevisit(),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	// Scrape the download file name and vew version number from Go download page
	c.OnHTML("a.download.downloadBox", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), fmt.Sprintf("%s.%s", osCpuType, extension)) {
			name, found := strings.CutPrefix(e.Attr("href"), "/dl/")
			if !found {
				log.Fatalln("Something went wrong getting the download file.")
			}
			dlFileName = name
			i := strings.Index(name, osCpuType)
			if i < 0 {
				log.Fatalln("Something went wrong getting the new version number.")
			}
			newVersion = name[2 : i-1]
		}
	})

	c.Visit("https://go.dev/dl/")

	c.OnHTML("tr.highlight", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, dlFileName) {
			e.ForEach("td", func(idx int, td *colly.HTMLElement) {
				if idx == 5 {
					dlFileCheckSum = td.Text
				}
			})
		}
	})

	c.Visit("https://go.dev/dl/")
}

func updateGo() {
	client := grab.NewClient()
	getFile := fmt.Sprintf("https://go.dev/dl/%s", dlFileName)
	client.UserAgent = "Mozilla/5.0"
	req, err := grab.NewRequest(os.TempDir(), getFile)
	if err != nil {
		panic(err)
	}

	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		panic(err)
	}

	defer os.Remove(resp.Filename)
	sha256Chksum := calculateSHA256(resp.Filename)
	if dlFileCheckSum != sha256Chksum {
		log.Fatalf("File validation failed!\n  Original checksum: %s\nCalculated checksum: %s\n", dlFileCheckSum, sha256Chksum)
	}
	parms := parameters{curVersion, newVersion, os.TempDir(), installDir, dlFileName, extension}
	done := false
	for i := 0; !done; i++ {
		maybeDone := false
		if i >= len(comments) {
			maybeDone = maybeDone || true
		} else {
			comment, err := template.New("comment").Parse(comments[i])
			if err != nil {
				log.Fatal(err)
			}
			err = comment.Execute(os.Stdout, parms)
			if err != nil {
				log.Fatal(err)
			}
		}
		if i >= len(commands) {
			maybeDone = maybeDone || true
		} else {
			command, err := template.New("command").Parse(commands[i])
			if err != nil {
				log.Fatal(err)
			}
			var cmdToRun strings.Builder
			err = command.Execute(&cmdToRun, parms)
			if err != nil {
				log.Fatal(err)
			}
			cmdAndArgsToRun := strings.Split(cmdToRun.String(), separator)
			if len(cmdAndArgsToRun) < 1 {
				log.Fatal("Command to run is empty.")
			}
			cmd := exec.Command(cmdAndArgsToRun[0], cmdAndArgsToRun[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmdErr := cmd.Run()
			if cmdErr != nil {
				if cmd.ProcessState.ExitCode() == 1602 {
					fmt.Println("Installation cancelled by user.")
					break
				}
				log.Fatal(cmdErr)
			}
		}
		done = !maybeDone
	}
	getCurrentVersion()
	fmt.Printf("Installed version is now %s\n", curVersion)
	fmt.Println("Done")
}

func calculateSHA256(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))

}
