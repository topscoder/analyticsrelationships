package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"

func crash(message string, err error) {
	fmt.Print(string(colorRed) + "[ERROR] " + message + string(colorReset) + "\n")
	panic(err)
}

func info(message string, silent bool) {
	if !silent {
		fmt.Print("[-] " + message + "\n")
	}
}

func banner() {
	data := `
██╗   ██╗ █████╗       ██╗██████╗
██║   ██║██╔══██╗      ██║██╔══██╗
██║   ██║███████║█████╗██║██║  ██║
██║   ██║██╔══██║╚════╝██║██║  ██║
╚██████╔╝██║  ██║      ██║██████╔╝
 ╚═════╝ ╚═╝  ╚═╝      ╚═╝╚═════╝

██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗███████╗
██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║██╔════╝
██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║███████╗
██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║╚════██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║███████║
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝

`
	data += "\033[32m> \033[0mGet related (sub)domains by looking at Google Analytics IDs\n"
	data += "\033[32m> \033[0mBy @JosueEncinar, @topscoder\n"

	println(data)
}

type UaIdentifier struct {
	UaCode       string
	OriginDomain string
}

// getURLResponse fetches the response body of a given URL.
func getURLResponse(url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 3}
	res, err := client.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extractGoogleTagManager extracts the Google Tag Manager URL or UA from the given URL.
func extractGoogleTagManager(targetURL string) (bool, []UaIdentifier) {
	var resultTagManager []UaIdentifier
	response, err := getURLResponse(targetURL)
	if err != nil {
		return false, resultTagManager
	}

	pattern := regexp.MustCompile(`www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+`)
	data := pattern.FindStringSubmatch(response)
	if len(data) > 0 {
		resultTagManager = append(
			resultTagManager,
			UaIdentifier{
				UaCode:       "https://" + strings.Replace(data[0], "ns.html", "gtm.js", -1),
				OriginDomain: targetURL,
			},
		)
	} else {
		pattern = regexp.MustCompile("GTM-[A-Z0-9]+")
		data = pattern.FindStringSubmatch(response)
		if len(data) > 0 {
			resultTagManager = append(
				resultTagManager,
				UaIdentifier{
					UaCode:       "https://www.googletagmanager.com/gtm.js?id=" + data[0],
					OriginDomain: targetURL,
				},
			)
		} else {
			pattern = regexp.MustCompile(`UA-\d+-\d+`)
			aux := pattern.FindAllStringSubmatch(response, -1)
			var result []UaIdentifier
			for _, r := range aux {
				result = append(
					result,
					UaIdentifier{
						UaCode:       r[0],
						OriginDomain: targetURL,
					},
				)
			}
			return true, result
		}
	}

	return false, resultTagManager
}

func getUA(urlGoogleTagManager string, urlOrigin string) []UaIdentifier {
	pattern := regexp.MustCompile("UA-[0-9]+-[0-9]+")
	response, _ := getURLResponse(urlGoogleTagManager)
	var result []UaIdentifier
	if response != "" {
		aux := pattern.FindAllStringSubmatch(response, -1)
		for _, r := range aux {
			result = append(result, UaIdentifier{UaCode: r[0], OriginDomain: urlOrigin})
		}
	} else {
		result = nil
	}
	return result
}

func cleanRelationShips(domains [][]string) []string {
	var allDomains []string
	for _, domain := range domains {
		allDomains = append(allDomains, strings.Replace(domain[0], "/relationships/", "", -1))
	}
	return allDomains
}

func getDomainsFromBuiltWith(id string) []string {
	pattern := regexp.MustCompile(`/relationships/[a-z0-9\-\_\.]+\.[a-z]+`)
	url := "https://builtwith.com/relationships/tag/" + id
	response, _ := getURLResponse(url)
	var allDomains []string = nil
	if response != "" {
		allDomains = cleanRelationShips(pattern.FindAllStringSubmatch(response, -1))
	}
	return allDomains
}

func getDomainsFromHackerTarget(id string) []string {
	url := "https://api.hackertarget.com/analyticslookup/?q=" + id
	response, _ := getURLResponse(url)
	var allDomains []string = nil
	if response != "" && !strings.Contains(response, "API count exceeded") {
		allDomains = strings.Split(response, "\n")
	}
	return allDomains
}

func getDomains(id string) []string {
	var allDomains []string = getDomainsFromBuiltWith(id)
	domains2 := getDomainsFromHackerTarget(id)
	if len(domains2) != 0 {
		for _, domain := range domains2 {
			if !contains(allDomains, domain) {
				allDomains = append(allDomains, domain)
			}
		}
	}
	return allDomains
}

func contains(data []string, value string) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func showDomains(ua string, chainMode bool) {
	allDomains := getDomains(ua)
	if len(allDomains) == 0 {
		// fmt.Println("No domains found for ua " + ua)
		return
	}

	for _, domain := range allDomains {
		if domain == "error getting results" {
			// crash("Server-side error on builtwith.com: error getting results", nil)
			continue
		}

		fmt.Println(domain + "," + ua)
	}
}

func start(url string, silent bool) {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	info("Analyzing url: "+url, silent)

	isUaResult, resultTagManager := extractGoogleTagManager(url)

	if len(resultTagManager) > 0 {
		var visited []string
		var allUAs []UaIdentifier

		if !isUaResult {
			urlGoogleTagManager := resultTagManager[0]

			info("URL with UA: "+urlGoogleTagManager.UaCode, silent)

			allUAs = getUA(urlGoogleTagManager.UaCode, url)
		} else {
			info("Found UA directly", silent)

			allUAs = resultTagManager
		}

		info("Obtaining information from builtwith and hackertarget\n", silent)

		for _, uaIdentifier := range allUAs {
			baseUA := strings.Join(strings.Split(uaIdentifier.UaCode, "-")[0:2], "-")

			if !contains(visited, baseUA) {
				visited = append(visited, baseUA)
				showDomains(baseUA, silent)
			}
		}

	} else {
		info("Tagmanager URL not found", silent)
	}
}

func main() {
	url := flag.String("url", "", "URL to extract Google Analytics ID from.")
	silent := flag.Bool("silent", false, "Don't print shizzle. Only what matters.")
	flag.Parse()

	if !*silent {
		banner()
	}

	if *url != "" {
		start(*url, *silent)
	} else {
		stat, _ := os.Stdin.Stat()

		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					crash("bufio couldn't read stdin correctly.", err)
				} else {
					start(scanner.Text(), *silent)
				}
			}
		}
	}
}
