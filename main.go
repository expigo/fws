package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
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

	formData := url.Values{
		"j_username": {*u},
		"j_password": {*p},
	}

	// formData.Add("_login_redirect_url", "https://www.filmweb.pl/user/"+*u)
	// formData.Add("_login_redirect_url", "https://www.filmweb.pl/film/Nietykalni-2011-583390")
	// formData.Add("_login_redirect_url", "https://www.filmweb.pl/Skazani.Na.Shawshank")
	formData.Add("_login_redirect_url", "https://www.filmweb.pl/ranking/film")

	req, err := http.NewRequest("POST", "https://www.filmweb.pl/j_login", strings.NewReader(formData.Encode()))
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

	var ratings data

	doc.Find(".filmPoster__filmLink").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 100 {
			return false
		}
		href, exists := s.Attr("href")
		if exists {
			var mi = movieInfo{URL: href}
			parse(href, &client, &mi)
			ratings = append(ratings, mi)
		}

		return true
	})

	result, err := json.MarshalIndent(ratings, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("test.json", result, 0644)

	fmt.Print("Done 👍")

}

func parse(url string, c *http.Client, mi *movieInfo) {

	resp, _ := c.Get("https://filmweb.pl" + url)
	if resp.StatusCode == http.StatusOK {

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		//  , ],{avg: 9.17, count: 43, limit: 4}] })
		userRating := regexp.MustCompile(`],{avg: (.*?), count: (.*?),`)
		userRatingRaw := userRating.FindStringSubmatch(bodyString)
		(*mi).Friends.Rating, _ = strconv.ParseFloat(userRatingRaw[1], 64)
		(*mi).Friends.Count, _ = strconv.Atoi(userRatingRaw[2])

		communityRating := regexp.MustCompile(`communityRateInfo:"(.*?)",communityRatingCountInfo:"(.*?) ocen[y]?"`)
		communityRatingRaw := communityRating.FindStringSubmatch(bodyString)
		(*mi).Community.Rating, _ = strconv.ParseFloat(strings.Replace(communityRatingRaw[1], ",", ".", -1), 64)
		(*mi).Community.Count, _ = strconv.Atoi(strings.Join(strings.Fields(communityRatingRaw[2]), ""))
	} else {
		panic(":O")
	}
}

type movieInfo struct {
	URL       string
	Friends   friendsRating
	Community communityRating
}

type friendsRating struct {
	Rating float64
	Count  int
}

type communityRating struct {
	Rating float64
	Count  int
}

type data []movieInfo
