package files

import (
	"bufio"
	"os"
	"strings"
)

type PropertiesFile struct {
	path    string
	entries map[string]string
}

func NewPropertiesFile(path string) *PropertiesFile {
	var file = new(PropertiesFile)

	file.path = path
	file.entries = make(map[string]string)

	return file
}

func (this *PropertiesFile) Put(key string, value string) {
	this.entries[key] = value
}

func (this *PropertiesFile) Get(key string) string {
	var value, found = this.entries[key]

	if found {
		return value
	} else {
		return ""
	}
}

func (this *PropertiesFile) Exists() bool {
	var _, openErr = os.Stat(this.path)

	if openErr == nil {
		return true
	} else {
		return false
	}
}

func (this *PropertiesFile) Save() error {
	var file, openError = os.Create(this.path)

	if openError != nil {
		return openError
	}

	defer file.Close()

	for k, v := range this.entries {
		file.WriteString(k + "=" + v + "\n")
	}

	return nil
}

func (this *PropertiesFile) Load() error {
	var file, openError = os.Open(this.path)

	if openError != nil {
		return openError
	}

	var scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		var line = scanner.Text()
		var index = strings.Index(line, "=")

		if index != -1 {
			var key = line[:index]
			var value = line[index+1:]

			this.entries[key] = value
		}
	}

	return nil
}
