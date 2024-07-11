package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const modDir string = "test-server/mods"

func (mod ModrinthMod) Install() (err error) {
	modIndexDir := filepath.Join(modDir, ".index")
	err = os.MkdirAll(modIndexDir, os.ModePerm)

	res, err := http.Get(mod.LatestVersion.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(filepath.Join(modDir, mod.LatestVersion.Filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	mod.writeIndex()

	return nil
}

func (mod ModrinthMod) writeIndex() {
	modIndexDir := filepath.Join(modDir, ".index")
	err := os.MkdirAll(modIndexDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}
	data := fmt.Sprintf(`filename = '%s'
name = '%s'
side = '%s'

[download]
hash = '%s'
hash-format = '%s'
mode = 'url'
url = '%s'

[update.modrinth]
mod-id = '%s'
version = '%s'`, mod.LatestVersion.Filename, mod.Title, mod.Side, mod.LatestVersion.Hash, mod.LatestVersion.HashAlgo, mod.LatestVersion.Url, mod.Id, mod.LatestVersion.Id)
	err = os.WriteFile(filepath.Join(modIndexDir, mod.Slug+".pw.toml"), []byte(data), os.ModePerm)

	if err != nil {
		log.Fatal(err)
		return
	}
}
