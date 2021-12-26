package altsrc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"

	"gopkg.in/yaml.v2"
)

type yamlSourceContext struct {
	FilePath string
}

// NewYamlSourceFromFile creates a new Yaml InputSourceContext from a filepath.
func NewYamlSourceFromFile(file string) (InputSourceContext, error) {
	ysc := &yamlSourceContext{FilePath: file}
	var results map[interface{}]interface{}
	err := readCommandYaml(ysc.FilePath, &results)
	if err != nil {
		return nil, fmt.Errorf("Unable to load Yaml file '%s': inner error: \n'%v'", ysc.FilePath, err.Error())
	}

	return &MapInputSource{file: file, valueMap: results}, nil
}

// NewYamlSourceFromFlagFunc creates a new Yaml InputSourceContext from a provided flag name and source context.
func NewYamlSourceFromFlagFunc(flagFileName string) func(context *cli.Context) (InputSourceContext, error) {
	return func(context *cli.Context) (InputSourceContext, error) {
		if context.IsSet(flagFileName) {
			filePath := context.String(flagFileName)
			return NewYamlSourceFromFile(filePath)
		}

		return defaultInputSource()
	}
}

func readCommandYaml(filePath string, container interface{}) (err error) {
	b, err := loadDataFrom(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, container)
	if err != nil {
		return err
	}

	err = nil
	return
}
