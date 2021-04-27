package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/spf13/cobra"
)

var templateString = `
package internal

var Data map[string][]byte

func init(){
	Data = map[string][]byte{
		{{range .}}
		"{{.Path}}":[]byte{{.Load}},
		{{end}}
	}
}

`

type Payload struct {
	Path string
	Load string
}

var Data []Payload = []Payload{}
var dirs []string = []string{}

var command = &cobra.Command{
	Short: "File encoding generator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dirs)
		tpl, err := template.New("Encoded").Parse(templateString)
		if err != nil {
			fmt.Println("Failed to create and parse template ", err)
		}
		wg := new(sync.WaitGroup)
		for _, dir := range dirs {
			wg.Add(1)
			go walkDir(dir, wg)
		}
		wg.Wait()
		file, err := os.Create("internal/generated.go")
		if err != nil {
			fmt.Sprintf("Failed to open file, %s", err)
			return
		}
		defer file.Close()
		err = tpl.Execute(file, Data)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	command.PersistentFlags().StringSliceVarP(&dirs, "dirs", "d", []string{"templates", "static"}, "Directories to walk through")
}

// Function to walk through all the files through directory
// and generate byte sequence to encode
func walkDir(dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		var encoded []string = []string{}
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Println("reading for file ", info.Name())
			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file %s", err)
				return err
			}
			for _, i := range data {
				encoded = append(encoded, fmt.Sprintf("%d", int(i)))
			}
			Data = append(Data, Payload{Path: info.Name(), Load: fmt.Sprintf("{%s}", strings.Join(encoded, ","))})
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directory ", dir)
	}
}

func main() {
	err := command.Execute()
	if err != nil {
		fmt.Println("Error Executing command ", err)
	}
}
