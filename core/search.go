package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type Cache struct {
	SkypeDir string `json:"skype_dir"`
}

func searchForSkypeDbs() []string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	baseDir := usr.HomeDir

	//Try to get path to Skype Dir from cache file
	cache := new(Cache)
	b, _ := ioutil.ReadFile(GetRootDir() + "/cache.json")
	json.Unmarshal(b, &cache)

	if cache.SkypeDir != "" {
		baseDir = cache.SkypeDir
	}

	var files []string

	filepath.Walk(baseDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() && strings.Contains(path, "Skype") {
			if f.Name() == "main.db" {
				files = append(files, path)
			}
		}
		return nil
	})

	//Cache path to Skype dir for quick search
	if cache.SkypeDir == "" && len(files) > 0 {
		cacheFile, err := CreateFile(GetRootDir() + "/cache.json")
		defer cacheFile.Close()

		if err == nil {
			fullPath := files[0]
			re := regexp.MustCompile("Skype/(.*?)/main.db")
			rm := re.FindStringSubmatch(fullPath)

			cache.SkypeDir = strings.Replace(filepath.Dir(fullPath), rm[1], "", 1)

			b, err := json.Marshal(cache)
			if err == nil {
				cacheFile.Write(b)
			}
		}
	}

	return files
}
