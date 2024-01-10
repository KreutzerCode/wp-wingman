package wordpressFinder

import (
	"io"
	"net/http"
	"regexp"
	"wp-wingman/utils"
)

func IsWordpressSite(url string, useRandomUserAgent bool) bool {
	isFound := utils.TestUrlForAvailability(url + "/wp-login.php", useRandomUserAgent)

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
