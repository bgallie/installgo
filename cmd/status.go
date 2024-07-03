/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var (
	curVersion     string
	newVersion     string
	dlFileName     string
	dlFileCheckSum string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("status called")
		getCurrentVersion()
		fmt.Printf("Current Version: %s\n", curVersion)
		scrapeLatestVersion()
		if curVersion != newVersion {
			fmt.Printf("A new version, %s, is available.\n", newVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func isCacheValid(cacheFile string) bool {
	fCache, err := os.Open(cacheFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fCache.Close()
	cacheInfo, err := fCache.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return time.Since(cacheInfo.ModTime()).Hours() <= maxCacheTime
}

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
