package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Index struct {
	Path    string
	Entries map[string]Entry
}

func LoadIndex(path string) (Index, error) {
	file, err := os.Open(path)
	if err != nil {
		// return Index{"", nil}, err
		GenerateIndex(path)
		file, err = os.Open(path)
		if err != nil {
			return Index{"", nil}, err
		}
	}
	defer file.Close()

	var entries map[string]Entry
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &entries)

	index := Index{path, entries}
	index.VerifyKeys()
	return index, nil
}

func (index Index) Save() error {
	json, err := json.Marshal(index.Entries)
	ioutil.WriteFile(index.Path, json, 0644)
	return err
}

func GenerateIndex(path string) error {
	root, err := os.Open("./content")
	if err != nil {
		err := os.Mkdir("./content", 0755)
		if err != nil {
			return err
		}
		root, _ = os.Open("./content")
	}
	defer root.Close()

	dirs, err := root.Readdirnames(0)
	index := Index{path, make(map[string]Entry)}
	for _, dir := range dirs {
		entry, err := GenerateEntry(dir)
		if err != nil {
			return err
		}
		index.Entries[dir] = entry
	}
	index.Save()
	return nil
}

func (index Index) VerifyKeys() {
	for key, entry := range index.Entries {
		entry.Key = key
		index.Entries[key] = entry
	}
}

type Entry struct {
	Key      string   `json:"key"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	Cover    string   `json:"cover"`
	Episodes []Source `json:"episodes"`
}

func GenerateEntry(key string) (Entry, error) {
	root, err := os.Open("./content/" + key)
	if err != nil {
		return Entry{}, err
	}
	defer root.Close()

	files, err := root.Readdirnames(0)
	entry := Entry{key, key, []string{}, "", []Source{}}
	for _, file := range files {
		entry.GenerateSource(file)
	}
	return entry, nil
}

type Source struct {
	Title string `json:"title"`
	Path  string `json:"source"`
}

func (e *Entry) GenerateSource(file string) {
	if file[len(file)-4:] != ".mp4" || file[len(file)-5:] != ".webm" {
		return
	}
	source := Source{file[:len(file)-4], "./content/" + e.Key + "/" + file}
	e.Episodes = append(e.Episodes, source)
}
