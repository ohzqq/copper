package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var base string
var outDir = "."
var id string
var bURL string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copper",
	Short: "download pics from a coppermine gallery",
	Long: `example: 
	copper -u 'https://tylerhoechlin.org/pictures/thumbnails.php' -d testdata/ -i 583 -b "april 2nd wondercon 2022"`,
	Run: func(cmd *cobra.Command, args []string) {
		if base != "" {
			base = flect.New(base).Underscore().String()
		}
		outDir = flect.New(outDir).Underscore().String()
		outDir = filepath.Join(outDir, base)
		err := os.MkdirAll(outDir, 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		println(outDir)

		siteURL, err := url.Parse(bURL)
		if err != nil {
			log.Fatal(err)
		}
		siteURL.Path = filepath.Dir(siteURL.Path)
		siteURL.RawQuery = ""

		c := colly.NewCollector()
		count := 1
		c.OnHTML("img.thumbnail", func(e *colly.HTMLElement) {
			dir := filepath.Dir(e.Attr("src"))
			name := e.Attr("alt")
			dl := siteURL.JoinPath(dir, name)
			//path := filepath.Join(u.Path, dir, name)
			//dl := *u
			//dl.Path = path
			fmt.Printf("%v\n", dl.String())

			resp, err := http.Get(dl.String())
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			name = dlPath(name, count)
			dlPix(name, resp.Body)
			fmt.Printf("saving %v\n", name)
			count++
			//resp
		})

		v := make(url.Values)
		v.Set("album", id)
		gal, err := url.Parse(bURL)
		if err != nil {
			log.Fatal(err)
		}
		gal.RawQuery = v.Encode()
		c.Visit(gal.String())
	},
}

func dlPath(name string, count int) string {
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	if base != "" {
		name = fmt.Sprintf("%s-%02d", base, count)
	}
	name = flect.New(name).Underscore().String()
	name = filepath.Join(outDir, name)
	name = fmt.Sprintf("%s%s", name, ext)
	return name
}

func dlPix(name string, body io.Reader) {
	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, body)
	if err != nil {
		log.Fatal(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&base, "base", "b", "", "basename of files")
	rootCmd.Flags().StringVarP(&outDir, "dir", "d", ".", "output dir")
	rootCmd.Flags().StringVarP(&bURL, "url", "u", "", "url of site")
	rootCmd.Flags().StringVarP(&id, "id", "i", "", "album id")
}
