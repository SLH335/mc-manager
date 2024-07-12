package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

type ModrinthMod struct {
	Id            string
	Slug          string
	Title         string
	Side          string
	LatestVersion ModrinthModVersion
}

type ModrinthModVersion struct {
	Id           string
	Filename     string
	Url          string
	Hash         string
	HashAlgo     string
	Dependencies []string
}

const baseUrl string = "https://api.modrinth.com/v2/"

func GetModrinthModInformation(ids []string, mcVersion string, printProgress bool) (mods []ModrinthMod, err error) {
	queryParams := url.Values{}
	queryParams.Set("ids", fmt.Sprintf("[\"%s\"]", strings.Join(ids, "\",\"")))

	jsonData, err := modrinthRequest(http.MethodGet, "projects?"+queryParams.Encode())
	if err != nil {
		return mods, err
	}

	modCount := len(jsonData.GetArray())
	for i, modData := range jsonData.GetArray() {
		clientSide := string(modData.GetStringBytes("client_side"))
		isClientSide := clientSide == "required" || clientSide == "optional"
		serverSide := string(modData.GetStringBytes("server_side"))
		isServerSide := serverSide == "required" || serverSide == "optional"
		side := ""
		if isClientSide && isServerSide {
			side = "both"
		} else if isClientSide {
			side = "client"
		} else {
			side = "server"
		}

		mod := ModrinthMod{
			Id:    string(modData.GetStringBytes("id")),
			Slug:  string(modData.GetStringBytes("slug")),
			Title: string(modData.GetStringBytes("title")),
			Side:  side,
		}
		if printProgress {
			fmt.Printf("\033[2K\rLoading mod [%d/%d] %s", i+1, modCount, mod.Title)
		}

		mod.LatestVersion, err = mod.getLatestVersion(mcVersion)
		if err != nil {
			return mods, err
		}

		mods = append(mods, mod)
	}
	if printProgress {
		fmt.Print("\n")
	}
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].Title < mods[j].Title
	})

	return mods, nil
}

func (mod ModrinthMod) getLatestVersion(mcVersion string) (modVersion ModrinthModVersion, err error) {
	path := fmt.Sprintf(`project/%s/version`, mod.Id)
	queryParams := url.Values{}
	queryParams.Set("loaders", fmt.Sprintf("[\"%s\"]", "fabric"))
	queryParams.Set("game_versions", fmt.Sprintf("[\"%s\"]", mcVersion))

	jsonData, err := modrinthRequest(http.MethodGet, path+"?"+queryParams.Encode())
	if err != nil {
		return ModrinthModVersion{}, err
	}

	found := false
	for _, file := range jsonData.GetArray("0", "files") {
		if file.GetBool("primary") && strings.HasSuffix(string(file.GetStringBytes("filename")), ".jar") {
			modVersion = ModrinthModVersion{
				Id:       string(jsonData.GetStringBytes("0", "id")),
				Filename: string(file.GetStringBytes("filename")),
				Url:      string(file.GetStringBytes("url")),
				Hash:     string(file.GetStringBytes("hashes", "sha512")),
				HashAlgo: "sha512",
			}
			found = true
			break
		}
	}

	for _, dependency := range jsonData.GetArray("0", "dependencies") {
		if string(dependency.GetStringBytes("dependency_type")) == "required" {
			modVersion.Dependencies = append(modVersion.Dependencies, string(dependency.GetStringBytes("project_id")))
		}
	}

	if !found {
		return ModrinthModVersion{}, fmt.Errorf("mod %s has no versions for mc %s", mod.Title, mcVersion)
	}

	return modVersion, nil
}

func modrinthRequest(method, path string) (data *fastjson.Value, err error) {
	req, err := http.NewRequest(method, baseUrl+path, nil)
	if err != nil {
		return &fastjson.Value{}, err
	}

	req.Header.Add("User-Agent", "slh335/mc-mod-manager")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return &fastjson.Value{}, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &fastjson.Value{}, err
	}

	var parser fastjson.Parser
	data, err = parser.Parse(string(resBody))
	if err != nil {
		return &fastjson.Value{}, err
	}

	return data, nil
}
