package leetcode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/user"
	"path"
)

var err error

func init() {
	http.DefaultClient.Jar, err = cookiejar.New(nil)
	checkErr(err)
	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Referer", req.URL.String())
		fmt.Println(req.URL.String())
		if len(via) >= 3 {
			return errors.New("stopped after 3 redirects")
		}
		return nil
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getCsrfToken(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "csrftoken" {
			return cookie.Value
		}
	}
	return ""
}

func getCachePath(f string) string {
	dir, err := os.UserCacheDir()
	checkErr(err)
	u, err := user.Current()
	if err == nil && u.HomeDir != "" {
		dir = path.Join(u.HomeDir, ".leetcode")
	}
	return path.Join(dir, f)
}

func getFilePath(filename string) string {
	if dir := path.Dir(filename); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			checkErr(err)
		}
	}
	return filename
}

func filePutContents(filename string, data []byte) {
	filename = getFilePath(filename)
	err = ioutil.WriteFile(filename, data, 0644)
	checkErr(err)
}

func jsonEncode(v interface{}) []byte {
	data, err := json.Marshal(v)
	checkErr(err)
	dst := bytes.Buffer{}
	err = json.Indent(&dst, data, "", "\t")
	checkErr(err)
	return dst.Bytes()
}
func saveCookies(cookies []*http.Cookie) {
	filePutContents(getCachePath(cookiesFile), jsonEncode(cookies))
}

func getCookies() (cookies []*http.Cookie) {
	b, err := ioutil.ReadFile(getCachePath(cookiesFile))
	checkErr(err)
	err = json.Unmarshal(b, &cookies)
	checkErr(err)
	return
}

func saveCredential(data url.Values) {
	filePutContents(getCachePath(credentialsFile), jsonEncode(data))
}

func getCredential() (data url.Values) {
	b, err := ioutil.ReadFile(getCachePath(credentialsFile))
	checkErr(err)
	err = json.Unmarshal(b, &data)
	checkErr(err)
	return
}
