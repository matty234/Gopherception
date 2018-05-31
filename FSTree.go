package main

import (
	"path/filepath"
	"io/ioutil"
	"log"
	"github.com/pkg/errors"
	"os"
)

const (
	File = "0"
	Directory = "1"
	Gif = "g"
)

type file struct {
	name string
	ext string
	size int64
}

type directory struct {
	path		string
	name 		string
	files 		[]file
	directories []directory
}

func (d *directory) iterate(i func(string, string)) {
	for _, f := range d.directories {
		i(f.name, Directory)
	}
	for _, f := range d.files {
		if f.ext == ".gif" {
			i(f.name, Gif)
		} else {
			i(f.name, File)
		}
	}
}

func (d *directory) addDirectory(addDirectory *directory) {
	d.directories = append(d.directories, *addDirectory)
}

func (d *directory) getIndex(workingDirectory string) (v []byte, err error) {
	for _, f := range d.files  {
		if f.name == "index.txt" {
			return ioutil.ReadFile(workingDirectory+d.path+"/"+f.name)
		}
	}
	return nil, errors.New("No index available")
}

func (d *directory) addFile(addFile *file) {
	d.files = append(d.files, *addFile)
}

func GetDirectoryAtPath(workingDirectory string, path string) *directory {
	stat, err := os.Stat(workingDirectory+path)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if stat.IsDir() {
		fileInfo, err := ioutil.ReadDir(workingDirectory+path)
		var returnDirectory = directory{
			path: path,
			name: ".",
		}
		if err != nil {
			log.Fatal(err)
			return &returnDirectory
		}
		for _, f := range fileInfo {
			if f.IsDir() {
				returnDirectory.addDirectory(&directory{
					name: f.Name(),
					path: path,
				})
			} else {
				returnDirectory.addFile(&file{
					name: f.Name(),
					size: f.Size(),
					ext: filepath.Ext(f.Name()),
				})
			}
		}
		return &returnDirectory
	}
	return nil
}