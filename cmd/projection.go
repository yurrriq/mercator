// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ghodss/yaml"
	dat "github.com/shanesiebken/mercator/data"
	"github.com/spf13/cobra"
)

// Projection represents the structure and replacements for
// generated values files for some chart
type Projection map[string]interface{}

// projectionCmd represents the projection command
var projectionCmd = &cobra.Command{
	Use:   "projection",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pf, err := cmd.Flags().GetString("projection")
		if err != nil {
			panic(err)
		}
		vf, err := cmd.Flags().GetString("values")
		if err != nil {
			panic(err)
		}
		d, err := cmd.Flags().GetString("destination")
		if err != nil {
			panic(err)
		}
		o, err := cmd.Flags().GetBool("overwrite")

		if err != nil {
			panic(err)
		}
		createProjection(pf, vf, d, o)

		r, err := cmd.Flags().GetBool("readme")
		if r {
			createReadme(pf, d)
		}
	},
}

func init() {
	rootCmd.AddCommand(projectionCmd)

	projectionCmd.Flags().StringP("projection", "f", "projection.yaml", "Projection file for templating")
	projectionCmd.Flags().StringP("values", "v", "src.values.yaml", "Source values file for templating")
	projectionCmd.Flags().StringP("destination", "d", "", "Destination for templated values files")
	projectionCmd.Flags().BoolP("overwrite", "o", true, "Overwrite existing destination files")

	projectionCmd.Flags().BoolP("readme", "r", true, `Write a (or overwrite an existing) Readme which \
includes deployment instructions for projections`)
}

// ReadProjection will parse YAML byte data into a Values.
func ReadProjection(data []byte) (proj Projection, err error) {
	err = yaml.Unmarshal(data, &proj)
	if len(proj) == 0 {
		proj = Projection{}
	}
	return
}

// ReadProjectionFile will parse a YAML file into a map of values.
func ReadProjectionFile(filename string) (Projection, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return ReadProjection(data)
}

func createProjection(pf, vf, d string, o bool) {
	p, err := ReadProjectionFile(pf)
	if err != nil {
		panic(err)
	}
	projs := p["projections"].([]interface{})
	files := make([]struct {
		path    string
		content []byte
	}, len(projs))
	count := 0
	for _, k := range projs {
		if k.(map[string]interface{})["name"] != nil {

			files[count] = struct {
				path    string
				content []byte
			}{
				path:    d + k.(map[string]interface{})["name"].(string) + ".values.yaml",
				content: templateProjection(vf, k.(map[string]interface{})),
			}
		}
		count = count + 1
	}

	for _, file := range files {
		if !o {
			if _, err := os.Stat(file.path); err == nil {
				// File exists and is okay. Skip it.
				continue
			}
		}
		if err := ioutil.WriteFile(file.path, file.content, 0644); err != nil {
			panic(err)
		}
	}
}

func templateProjection(vf string, data interface{}) []byte {
	top := map[string]interface{}{
		"Projection": data,
	}
	var tpl bytes.Buffer
	templ, err := template.New("values").Delims("[[", "]]").ParseFiles(vf)
	if err != nil {
		panic(err)
	}
	_, fn := filepath.Split(vf)
	err = templ.ExecuteTemplate(&tpl, fn, top)
	if err != nil {
		panic(err)
	}
	return tpl.Bytes()
}

func createReadme(pf, d string) {
	p, err := ReadProjectionFile(pf)
	if err != nil {
		panic(err)
	}

	file := struct {
		path    string
		content []byte
	}{
		path:    d + "README.adoc",
		content: templateReadme(d, p),
	}

	if err := ioutil.WriteFile(file.path, file.content, 0644); err != nil {
		panic(err)
	}
}

func templateReadme(d string, data interface{}) []byte {

	top := map[string]interface{}{
		"Readme": data,
	}

	top["Readme"].(Projection)["destination"] = "./"
	if d != "" {
		top["Readme"].(Projection)["destination"] = d
	}
	r, err := dat.Asset("templates/README.template")
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	templ, err := template.New("readme").Parse(string(r))
	if err != nil {
		panic(err)
	}
	err = templ.Execute(&tpl, top)
	if err != nil {
		panic(err)
	}
	return tpl.Bytes()
}
