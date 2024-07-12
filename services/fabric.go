package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/valyala/fastjson"
)

func DownloadLatestFabricServer(mcVersion, dir string) (filename string, err error) {
	serverUrl, err := getLatestFabricServerUrl(mcVersion)
	if err != nil {
		return "", err
	}
	res, err := http.Get(serverUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	filename = strings.Split(res.Header["Content-Disposition"][0], "\"")[1]

	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	writeStartScript(filename, dir)

	return filename, nil
}

func UpdateFabricServer(mcVersion string) (err error) {
	loaderVersion, err := getLatestFabricVersion("loader")
	if err != nil {
		return err
	}
	installerVersion, err := getLatestFabricVersion("installer")
	if err != nil {
		return err
	}
	latestFilename := fmt.Sprintf("fabric-server-mc.%s-loader.%s-launcher.%s.jar", mcVersion, loaderVersion, installerVersion)

	filenames, err := filepath.Glob(fmt.Sprintf("fabric-server-mc.%s-loader.*-launcher.*.jar", mcVersion))
	if err != nil {
		return err
	}
	slices.Sort(filenames)
	currentFilename := filenames[len(filenames)-1]

	if latestFilename > currentFilename {
		_, err = DownloadLatestFabricServer("1.21", "./")
		if err != nil {
			return err
		}
		os.Remove(currentFilename)
		fmt.Println("Updated server software to", latestFilename)
	} else {
		fmt.Println("Server software is up to date")
	}

	return nil
}

func writeStartScript(filename, dir string) {
	os.WriteFile(filepath.Join(dir, "start.sh"), []byte(fmt.Sprintf("java -Xmx2G -jar %s nogui", filename)), 0755)
}

func getLatestFabricServerUrl(mcVersion string) (url string, err error) {
	loaderVersion, err := getLatestFabricVersion("loader")
	if err != nil {
		return "", err
	}
	installerVersion, err := getLatestFabricVersion("installer")
	if err != nil {
		return "", err
	}
	url = fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", mcVersion, loaderVersion, installerVersion)
	return url, nil
}

func getLatestFabricVersion(property string) (version string, err error) {
	res, err := http.Get("https://meta.fabricmc.net/v2/versions/" + property)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var parser fastjson.Parser
	jsonData, err := parser.Parse(string(resBody))
	if err != nil {
		return "", err
	}

	version = string(jsonData.GetStringBytes("0", "version"))
	return version, nil
}
