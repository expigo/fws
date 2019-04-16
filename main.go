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
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

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
	data.Add("_login_redirect_url", "https://www.filmweb.pl/film/Nietykalni-2011-583390")
	// data.Add("_login_redirect_url", "https://www.filmweb.pl/Skazani.Na.Shawshank")

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

	doc, _ := goquery.NewDocumentFromResponse(resp)
	htmlContent, _ := doc.Html()
	io.Copy(out, strings.NewReader(htmlContent))

	//  , ],{avg: 9.17, count: 43, limit: 4}] })
	userRating := regexp.MustCompile(`],{avg: (.*?), count: (.*?),`)
	userRatingRaw := userRating.FindStringSubmatch(htmlContent)
	fmt.Println(userRatingRaw[1])
	fmt.Println(userRatingRaw[2])

	communityRating := regexp.MustCompile(`communityRateInfo:"(.*?)",communityRatingCountInfo:"(.*?) ocen."`)
	communityRatingRaw := communityRating.FindStringSubmatch(htmlContent)
	fmt.Println(communityRatingRaw[1])
	fmt.Println(communityRatingRaw[2])

	fmt.Print("Done 👍")

}
