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
var page string
var dry bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copper",
	Short: "download pics from a coppermine gallery",
	Long: `example: 
	copper -u 'https://tylerhoechlin.org/pictures/ -d testdata/ -i 583 -n "april 2nd wondercon 2022"`,
	Run: func(cmd *cobra.Command, args []string) {
		baseURL, err := url.Parse(bURL)
		if err != nil {
			log.Fatal(err)
		}

		//basePath := filepath.Dir(baseURL.Path)
		//fmt.Printf("base path %s\n", basePath)

		albumURL := baseURL.JoinPath("thumbnails.php")
		v := make(url.Values)
		v.Set("album", id)
		if page != "" {
			v.Set("page", page)
		}
		albumURL.RawQuery = v.Encode()
		fmt.Printf("album url %s\n", albumURL)

		if base != "" {
			base = flect.New(base).Underscore().String()
		}

		outDir = flect.New(outDir).Underscore().String()
		outDir = filepath.Join(outDir, base)
		err = os.MkdirAll(outDir, 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		fmt.Printf("output dir %s\n", outDir)

		//albumURL.Path = albumsPath

		c := colly.NewCollector()
		//count := 1
		c.OnHTML("img.thumbnail", func(e *colly.HTMLElement) {
			albumURL.RawQuery = ""
			dir := filepath.Dir(e.Attr("src"))
			name := e.Attr("alt")
			dl := baseURL.JoinPath(dir, name)
			fmt.Printf("dl url %s\n", dl)

			name = dlPath(name)
			fmt.Printf("saving %v\n", name)

			if !dry {
				resp, err := http.Get(dl.String())
				if err != nil {
					log.Print(err)
				}
				defer resp.Body.Close()
				dlPix(name, resp.Body)
			}
			//count++
			//resp
		})

		c.Visit(albumURL.String())
	},
}

func dlPath(name string) string {
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
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
	rootCmd.Flags().StringVarP(&base, "name", "n", "", "basename of files")
	rootCmd.Flags().StringVarP(&outDir, "dir", "d", ".", "output dir")
	rootCmd.Flags().StringVarP(&bURL, "url", "u", "", "url of site")
	rootCmd.Flags().StringVarP(&id, "id", "i", "", "album id")
	rootCmd.Flags().StringVarP(&page, "page", "p", "", "page num")
	rootCmd.Flags().BoolVar(&dry, "dry-run", false, "do a dry run (no files will be downloaded")
}
