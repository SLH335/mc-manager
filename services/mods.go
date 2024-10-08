package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func InitIndex(loader, mcVersion string) (err error) {
	modIndexDir, err := getModIndexDir()
	if err != nil {
		return err
	}

	data := []byte(fmt.Sprintf("loader = '%s'\nmc-version = '%s'", loader, mcVersion))
	err = os.WriteFile(filepath.Join(modIndexDir, "mc-manager.index.toml"), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func GetVersion() (loader, mcVersion string, err error) {
	loader, mcVersion, err = getVersionFromIndex()
	if err != nil {
		loader, mcVersion, err = DetectVersion()
		if err != nil {
			return "", "", err
		}
		InitIndex(loader, mcVersion)
	}
	return loader, mcVersion, nil
}

func getVersionFromIndex() (loader, mcVersion string, err error) {
	modIndexDir, err := getModIndexDir()
	if err != nil {
		return "", "", err
	}

	data, err := os.ReadFile(filepath.Join(modIndexDir, "mc-manager.index.toml"))
	if err != nil {
		return "", "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		vals := strings.Split(line, " = ")
		if len(vals) > 1 {
			val := strings.TrimSpace(vals[1])
			val = val[1 : len(val)-1]
			switch vals[0] {
			case "loader":
				loader = val
			case "mc-version":
				mcVersion = val
			}
		}
	}

	if loader == "" || mcVersion == "" {
		return loader, mcVersion, fmt.Errorf("index does not contain version information")
	}

	return loader, mcVersion, nil
}

func DetectVersion() (loader, mcVersion string, err error) {
	files, _ := filepath.Glob("*fabric*.jar")
	if len(files) > 0 {
		slices.Sort(files)
		file := files[len(files)-1]

		pattern := regexp.MustCompile("mc\\.(1(\\.[0-9]{1,2}){1,2})")
		version := pattern.FindStringSubmatch(file)
		if version[0] != "" {
			return "fabric", version[1], nil
		}
	}

	return "", "", fmt.Errorf("failed to detect mod loader")
}

func (mod ModrinthMod) Install() (err error) {
	modDir, err := getModDir()
	if err != nil {
		return err
	}

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

func UpdateMods(modIds []string, loader, mcVersion string) (err error) {
	modDir, err := getModDir()
	if err != nil {
		return err
	}

	mods := []ModrinthMod{}
	if len(modIds) == 0 {
		mods, err = getFromIndex()
		if err != nil {
			return err
		}
	}

	updated := 0
	for i, mod := range mods {
		fmt.Printf("\033[2K\rChecking mod [%d/%d] %s", i+1, len(mods), mod.Title)
		latestVersion, err := mod.getLatestVersion(loader, mcVersion)
		if err != nil {
			return err
		}
		if mod.LatestVersion.Id != latestVersion.Id {
			fmt.Printf("Updating mod: %s -> %s\n", mod.LatestVersion.Filename, latestVersion.Filename)
			updated++
			oldFilename := mod.LatestVersion.Filename
			mod.LatestVersion = latestVersion
			mod.Install()
			if mod.LatestVersion.Filename != oldFilename {

				os.Remove(filepath.Join(modDir, oldFilename))
			}
		}
	}
	if len(mods) > 0 {
		fmt.Print("\n")
	}
	if updated == 0 {
		fmt.Println("All mods up to date")
	} else {
		fmt.Printf("Updated %d mods\n", updated)
	}

	return nil
}

func getModDir() (modDir string, err error) {
	wd, err := os.Getwd()
	if err == nil && filepath.Base(wd) == "mods" {
		return "./", nil
	}

	stat, err := os.Stat("mods")
	if err == nil && stat.IsDir() {
		return "./mods", nil
	}

	return "", fmt.Errorf("no valid mod directory found")
}

func getModIndexDir() (modIndexDir string, err error) {
	modDir, err := getModDir()
	if err != nil {
		return "", err
	}
	modIndexDir = filepath.Join(modDir, ".index")
	err = os.MkdirAll(modIndexDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return modIndexDir, nil
}

func getFromIndex() (mods []ModrinthMod, err error) {
	modIndexDir, err := getModIndexDir()
	if err != nil {
		return mods, err
	}

	files, err := os.ReadDir(modIndexDir)
	if err != nil || len(files) == 0 {
		return mods, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".pw.toml") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(modIndexDir, file.Name()))
		if err != nil {
			continue
		}
		mod := ModrinthMod{
			Slug: strings.ReplaceAll(file.Name(), ".pw.toml", ""),
		}
		for _, line := range strings.Split(string(data), "\n") {
			vals := strings.Split(line, " = ")
			if len(vals) > 1 {
				val := strings.TrimSpace(vals[1])
				val = val[1 : len(val)-1]
				switch vals[0] {
				case "mod-id":
					mod.Id = val
				case "name":
					mod.Title = val
				case "side":
					mod.Side = val
				case "version":
					mod.LatestVersion.Id = val
				case "filename":
					mod.LatestVersion.Filename = val
				case "hash":
					mod.LatestVersion.Hash = val
				case "hash-format":
					mod.LatestVersion.HashAlgo = val
				case "url":
					mod.LatestVersion.Url = val
				}
			}
		}
		mods = append(mods, mod)
	}
	return mods, nil
}

func (mod ModrinthMod) writeIndex() (err error) {
	modIndexDir, err := getModIndexDir()
	if err != nil {
		return err
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

	return err
}
