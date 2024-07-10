package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fastjson"
)

func DownloadLatestFabricServer(mcVersion, dir string) (fileName string, err error) {
	serverUrl, err := getLatestFabricServerUrl(mcVersion)
	if err != nil {
		return "", err
	}
	res, err := http.Get(serverUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	fileName = strings.Split(res.Header["Content-Disposition"][0], "\"")[1]

	file, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return fileName, nil
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
