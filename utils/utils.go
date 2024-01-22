package utils

import (
	"bufio"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
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

func DoesRemoteFileExist(url string, useRandomUserAgent bool) (bool, error) {
	client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, err
    }

    if useRandomUserAgent {
        rand.Seed(time.Now().Unix())
        randomUserAgent := userAgents[rand.Intn(len(userAgents))]
        req.Header.Set("User-Agent", randomUserAgent)
    }

    resp, err := client.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
        return true, nil
    } else {
        return false, nil
    }
}

func FetchReadme(url string, useRandomUserAgent bool) (interface{}, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
	if useRandomUserAgent {
		randomUserAgent := userAgents[rand.Intn(len(userAgents))]
		req.Header.Set("User-Agent", randomUserAgent)
	}

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return false, nil
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return string(body), nil
}

func TestUrlForAvailability(url string, useRandomUserAgent bool) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", url, nil)

	if useRandomUserAgent {
		randomUserAgent := userAgents[rand.Intn(len(userAgents))]
		req.Header.Set("User-Agent", randomUserAgent)
	}

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

func GetUserInputYesNo() bool {
    reader := bufio.NewReader(os.Stdin)
    answer, _ := reader.ReadString('\n')
    answer = strings.ToLower(strings.TrimSpace(answer))
    if strings.HasPrefix(answer, "y") {
        return true
    } else {
        return false
    }
}