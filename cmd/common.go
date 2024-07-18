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

	"github.com/gocolly/colly"
)

func unCache(URL string) {
	sum := sha1.Sum([]byte(URL))
	hash := hex.EncodeToString(sum[:])
	dir := path.Join(cacheDir, hash[:2])
	filename := path.Join(dir, hash)
	if isCacheValid(filename) {
		return
	}
	log.Println("Deleting cached file:", filename)
	if err := os.Remove(filename); err != nil {
		log.Fatal(err)
	}
}

func getCurrentVersion() {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		log.Println(err)
	}
	ver := strings.Split(string(out), " ")
	curVersion = strings.TrimPrefix(ver[2], "go")
	osCpuType = strings.TrimSuffix(ver[3], "\n")
	osCpuType = strings.ReplaceAll(osCpuType, "/", "-")
}

func scrapeLatestVersion() {
	unCache("https://go.dev/dl/")
	c := colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.AllowURLRevisit(),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	// c2 := c.Clone()

	c.OnHTML("div.toggleVisible", func(e *colly.HTMLElement) {
		nv, found := strings.CutPrefix(e.Attr("id"), "go")
		if newVersion == "" && found {
			newVersion = nv
			dlFileName = fmt.Sprintf("go%s.%s.tar.gz", newVersion, osCpuType)
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
