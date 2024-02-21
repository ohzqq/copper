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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copper",
	Short: "download pics from a coppermine gallery",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if outDir != "." {
			outDir = flect.New(outDir).Underscore().String()
			err := os.Mkdir(outDir, 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
		}
		println(outDir)
		u, err := url.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}
		u.Path = filepath.Dir(u.Path)
		u.RawQuery = ""

		c := colly.NewCollector()
		count := 1
		c.OnHTML("img.thumbnail", func(e *colly.HTMLElement) {
			dir := filepath.Dir(e.Attr("src"))
			name := e.Attr("alt")
			dl := u.JoinPath(dir, name)
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
		c.Visit(args[0])
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.copper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&base, "base", "b", "", "basename of files")
	rootCmd.Flags().StringVarP(&outDir, "dir", "d", ".", "output dir")
}
