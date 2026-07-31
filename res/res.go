package res

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

/*
#include "res.h"
#include "stdlib.h"
*/
import "C"

type resource struct {
	data interface{}
	err  error
}

var (
	resources map[string]resource = map[string]resource{}
)

const tmp = ".tmp/"

//go:embed */**
var resfs embed.FS

func Init() error {
	C.Init()

	err := LoadDir(".")
	if err != nil {
		log.Printf("%+v\n", err)
	}
	files, err := os.ReadDir("lib")
	if err != nil {
		log.Printf("%+v\n", err)
	}
	for _, entry := range files {
		err = ReadArchive("lib/" + entry.Name())
		if err != nil {
			log.Printf("%+v\n", err)
		}
	}

	return nil
}

func WriteFs() error {
	for path, res := range resources {
		if res.err == nil {
			err := os.MkdirAll(filepath.Dir(tmp+path), os.ModePerm)
			if err != nil {
				return err
			}

			err = os.WriteFile(tmp+path, res.data.([]byte), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CleanFs() {
	os.RemoveAll(tmp)
}

func LoadDir(cwd string) error {
	dir, err := resfs.ReadDir(cwd)
	if err != nil {
		return err
	}
	for _, file := range dir {
		path := file.Name()
		if cwd != "." {
			path = cwd + "/" + path
		}
		if file.IsDir() {
			LoadDir(path)
		} else {
			data, err := resfs.ReadFile(path)
			resources[path] = resource{data, err}
		}
	}
	return nil
}

func GetRes(path string) (interface{}, error) {
	if res, ok := resources[path]; ok {
		if res.err == nil && strings.Contains(path, ".obj") {
			return raylib.LoadModel(path), nil
		}
		if res.err == nil && strings.Contains(path, ".iqm") {
			return raylib.LoadModel(path), nil
		}
		if res.err == nil && strings.Contains(path, ".gltf") {
			return raylib.LoadModel(path), nil
		}
		if res.err == nil && strings.Contains(path, "glb") {
			return raylib.LoadModel(path), nil
		}
		if res.err == nil && strings.Contains(path, "vox") {
			return raylib.LoadModel(path), nil
		}
		if res.err == nil && strings.Contains(path, "m3d") {
			return raylib.LoadModel(path), nil
		}

		return res.data, res.err
	} else {
		return nil, fmt.Errorf("Resource not found: %s\n", path)
	}
}

func ReadArchive(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if strings.Contains(path, ".zip") {
		return ReadZip(data)
	}
	return fmt.Errorf("Archive not supported: %s\n", path)
}

func ReadZip(data []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		data, err := readZipFile(f)
		resources[f.Name] = resource{data, err}
	}
	return nil
}

func readZipFile(f *zip.File) (interface{}, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	b, readErr := io.ReadAll(rc)
	closeErr := rc.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return b, nil
}

//export GetData
func GetData(cfile *C.char, cdir *C.char) *C.char {
	file := C.GoString(cfile)
	dir := C.GoString(cdir)
	path := dir
	if strings.Compare("", dir) != 0 {
		path = path + "/"
	}
	path = path + file

	pos := strings.Index(path, "./")
	if pos == 0 {
		path = path[2:]
	}

	if res, ok := resources[path]; ok && res.err == nil {
		return C.CString(string(res.data.([]byte)))
	}

	return nil
}
