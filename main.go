package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	//"os"
)

func main() {
	_ = download_player_skin()

	fmt.Println("按回车键退出...")
	fmt.Scanln()
}

func download_player_skin() error {
	fmt.Println("输入玩家名称")
	var Name string
	fmt.Scanln(&Name)
	uuid, err := getid(Name)
	if err != nil {
		fmt.Printf("无法获取 UUID: %s\n", err)
		return err
	}
	value := getvalue(uuid)
	url := geturl(value)
	err = downlond(url, Name)
	if err == nil {
		fmt.Println("下载成功")
	} else {
		fmt.Println("下载失败")
		return err
	}

	return nil
}

func getid(name_in string) (string, error) {
	var name string = name_in
	//fmt.Scan(&Name)
	json_idname, err := http.Get(fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", name))
	if err != nil {
		return "", err
	}
	defer json_idname.Body.Close()
	fmt.Println(fmt.Sprintf("请求成功%d", json_idname.StatusCode))
	fmt.Println(fmt.Sprintf("请求地址%v", json_idname.Request.URL))
	//fmt.Println(json_idname.StatusCode)
	//fmt.Println(json_idname)
	body, _ := io.ReadAll(json_idname.Body)
	//fmt.Println(string(body))

	type id_json struct {
		Id   string `json_idname:"Id"`
		Name string `json_idname:"Name"`
	}
	var idnane id_json
	err = json.Unmarshal(body, &idnane)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("uuid %s", idnane.Id))
	fmt.Println(fmt.Sprintf("name %s", idnane.Name))
	return idnane.Id, nil
}

func getvalue(uuid string) string {
	json_value, err := http.Get(fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid))
	if err != nil {
		fmt.Println(err)
	}
	defer json_value.Body.Close()
	body, _ := io.ReadAll(json_value.Body)
	type value_json struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Properties []struct {
			Name      string `json:"name"`
			Value     string `json:"value"`
			Signature string `json:"signature"`
		} `json:"properties"`
	}
	var value value_json
	err = json.Unmarshal(body, &value)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(value.Properties[0].Value)
	return value.Properties[0].Value
}
func geturl(value string) string {
	type url_json struct {
		Timestamp   int64  `json:"timestamp"`
		ProfileId   string `json:"profileId"`
		ProfileName string `json:"profileName"`
		Textures    struct {
			SKIN struct {
				Url      string `json:"url"`
				Metadata struct {
					Model string `json:"model"`
				} `json:"metadata"`
			} `json:"SKIN"`
		} `json:"textures"`
	}
	var url url_json
	base64_json, _ := base64.StdEncoding.DecodeString(value)
	err := json.Unmarshal(base64_json, &url)
	if err != nil {
		fmt.Println(err)
	}
	return url.Textures.SKIN.Url
}
func downlond(url string, file_name string) error {
	file_re, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer file_re.Body.Close()
	file_by, err := io.ReadAll(file_re.Body)
	if err != nil {
		fmt.Println(err)
	}

	file_path := fmt.Sprintf("%s.png", file_name)
	err = os.WriteFile(file_path, file_by, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
