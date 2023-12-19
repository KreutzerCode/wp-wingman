package wordpressFinder

import (
	"io"
	"math/rand"
	"net/http"
	"regexp"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:57.0) Gecko/20100101 Firefox/57.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:57.0) Gecko/20100101 Firefox/57.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109",
}

func IsWordpressSite(url string) bool {
	isFound := TestUrlForAvailability(url + "/wp-login.php")

	if !isFound {
		//check default wordpress readme
		isFound, _ = checkContentOnStringEvidence(url+"/readme.html", `wp-admin/images/wordpress-logo\.png`)
	}

	if !isFound {
		//check wp-content path evidence
		isFound, _ = checkContentOnStringEvidence(url, url+`wp-content/`)
	}

	return isFound
}

func checkContentOnStringEvidence(url string, pattern string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Check for the HTML element using a regular expression
	re := regexp.MustCompile(pattern)
	found := re.Match(body)

	return found, nil
}

func TestUrlForAvailability(url string) bool {
	randomUserAgent := userAgents[rand.Intn(len(userAgents))]

	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", url, nil)
	req.Header.Set("User-Agent", randomUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}
