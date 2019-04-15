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

	resp, err := client.PostForm("https://www.filmweb.pl/j_login", url.Values{
		"j_username": {*u},
		"j_password": {*p},
	})
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
