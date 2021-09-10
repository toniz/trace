/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package trace

import (
	"io/ioutil"
	"os"
	"path"
    "errors"

    "strings"
    "encoding/json"
    "encoding/xml"
    "gopkg.in/yaml.v2"
)

// Read File And Parse To Struct
func (c *Otel) LoadFile(file string, l interface{}) error {
	f, err := c.Read(file)
	if err != nil {
		return err
	}

	err = c.Parse(path.Ext(file), f, l)
	return err
}

// Read a file, And return file content by []byte
func (c *Otel) Read(file string) ([]byte, error) {
	s, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	var content []byte
	if !s.IsDir() {
		fileHandler, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		defer fileHandler.Close()
		content, err = ioutil.ReadAll(fileHandler)
	}
	return content, err
}

// File Parse
func (c *Otel) Parse(name string, text []byte, l interface{}) error {
    switch name {
    case "json", ".json":
        t := string(text)
        t = strings.Replace(t, "\n", " ", -1)
        t = strings.Replace(t, "\r", " ", -1)
        err := json.Unmarshal([]byte(t), l)
        return err
    case "xml", ".xml":
        err := xml.Unmarshal(text, l)
        return err
    case "yaml", ".yaml":
        err := yaml.Unmarshal(text, l)
        return err
    default:
        return errors.New("No Decoder Being Add.")
    }
}

