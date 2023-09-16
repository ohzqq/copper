package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var base string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copper",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		u, err := url.Parse(args[1])
		if err != nil {
			log.Fatal(err)
		}
		u.Path = filepath.Dir(u.Path)
		u.RawQuery = ""

		c := colly.NewCollector()
		c.OnHTML("img.thumbnail", func(e *colly.HTMLElement) {
			dir := filepath.Dir(e.Attr("src"))
			name := e.Attr("alt")
			path := filepath.Join(u.Path, dir, name)
			dl := *u
			dl.Path = path
			fmt.Printf("%v\n", dl.String())

			resp, err := http.Get(dl.String())
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			n := args[0] + name
			out, err := os.Create(n)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("saving %v\n", n)
			//resp
		})
		c.Visit(args[1])
	},
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
	rootCmd.Flags().StringVarP(&base, "name", "n", "", "basename of files")
}
