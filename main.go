package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const baseURL = "https://filmweb.pl"

var u = flag.String("u", "", "uesrname")
var p = flag.String("p", "", "password")

func main() {
	flag.Parse()

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Jar: jar,
	}

	data := url.Values{
		"j_username": {*u},
		"j_password": {*p},
	}

	// data.Add("_login_redirect_url", "https://www.filmweb.pl/user/"+*u)
	data.Add("_login_redirect_url", "https://www.filmweb.pl/Skazani.Na.Shawshank")

	req, err := http.NewRequest("POST", "https://www.filmweb.pl/j_login", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "GOofy, curious but harmless bot 🐱‍🏍")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("resp.html")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	fmt.Print("Done 👍")

}
