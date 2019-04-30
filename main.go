package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const globalUsage = `
Check to see if there is an updated version available for installed charts.
`

var outputFile string
var inputFile string
var templateFile string
var version = "canary"
var commit string
var date string

// ChartEntry is an individual chart entry
type ChartEntry struct {
	APIVersion  string    `yaml:"apiVersion"`
	AppVersion  string    `yaml:"appVersion"`
	Created     time.Time `yaml:"created"`
	Description string    `yaml:"description"`
	Digest      string    `yaml:"digest"`
	Name        string    `yaml:"name"`
	Urls        []string  `yaml:"urls"`
	Version     string    `yaml:"version"`
}

// Charts is a record of all charts and a few metadata
type Charts struct {
	APIVersion string                  `yaml:"apiVersion"`
	Entries    map[string][]ChartEntry `yaml:"entries"`
	Generated  time.Time               `yaml:"generated"`
}

func main() {
	rootDir, _ := os.Getwd()

	cmd := &cobra.Command{
		Use:   "repo-html [flags]",
		Short: fmt.Sprintf("Generates an html file for a repo (helm-repo-html %s)", version),
		RunE:  run,
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print current helm-repo-html version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version:", version)
			fmt.Println("commit:", commit)
			fmt.Println("date:", date)
		},
	}
	cmd.AddCommand(versionCmd)

	cmd.Flags().StringVarP(&outputFile, "output", "o", path.Join(rootDir, "index.html"), "output filename")
	cmd.Flags().StringVarP(&inputFile, "input", "i", path.Join(rootDir, "index.yaml"), "input filename")
	cmd.Flags().StringVarP(&templateFile, "template", "t", path.Join(rootDir, "index.tpl"), "template (go html/template format) filename")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	var err error
	var htmlTemplate *template.Template
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		log.Info("Template file not found, using default")
		htmlTemplate = defaultHTMLTemplate
	} else {
		log.Info("Template file: " + templateFile)
		htmlTemplate = template.Must(template.New(path.Base(templateFile)).ParseFiles(templateFile))
	}

	log.Info("Reading " + inputFile)
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Error reading yaml file: %v", err)
	}

	charts := Charts{}
	log.Info("Processing yaml")
	err = yaml.Unmarshal([]byte(data), &charts)
	if err != nil {
		return fmt.Errorf("Error processing yaml file: %v", err)
	}
	var outputHandle *os.File

	if outputFile == "-" {
		outputHandle = os.Stdout
		log.Info("Outputting template to stdout")
	} else {
		log.Info("Creating " + outputFile)
		outputHandle, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("Error creating index file: %v", err)
		}
		log.Info("Outputting template " + outputFile)
	}

	err = htmlTemplate.Execute(outputHandle, charts)
	if err != nil {
		return fmt.Errorf("Error creating template: %v", err)
	}
	return err
}

var (
	defaultHTMLTemplate = template.Must(template.New("htmlTemplate").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <title>Helm Charts</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.10.0/github-markdown.min.css" />
    <style>
      .markdown-body {
        box-sizing: border-box;
        min-width: 200px;
        max-width: 980px;
        margin: 0 auto;
        padding: 45px;
      }
      @media (max-width: 767px) {
        .markdown-body {
          padding: 15px;
        }
      }
    </style>
  </head>

  <body>

    <section class="markdown-body">
      <h1>Helm Charts</h1>

      <h2>Usage</h2>
      <pre lang="no-highlight"><code>
        helm repo add halkeye https://halkeye.github.io/helm-charts
      </code></pre>

      <p>These are presented as is. Anyone is free to use them, and make suggestions, but they were created for my own use. At some point I want to submit them to the actual helm charts repo.</p>

      <h2>Charts</h2>

      <ul>
			{{range $key, $chartEntry := .Entries }}
				<li>
					<p>
						{{ (index $chartEntry 0).Name }}
						(<a href="{{ (index (index $chartEntry 0).Urls 0) }}" title="{{ (index (index $chartEntry 0).Urls 0) }}">
						{{ (index $chartEntry 0).Version }}
						@
						{{ (index $chartEntry 0).AppVersion }}
						</a>)
					</p>
					<p>
						{{ (index $chartEntry 0).Description }}
					</p>
				</li>
			{{end}}
      </ul>
    </section>
		<time datetime="{{ .Generated.Format "2006-01-02T15:04:05" }}" pubdate id="generated">{{ .Generated.Format "Mon Jan 2 2006 03:04:05PM MST-07:00" }}</time>
  </body>
</html>
`))
)
