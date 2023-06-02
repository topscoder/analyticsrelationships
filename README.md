# AnalyticsRelationships

<p align="center">
  <a href="https://golang.org/dl/#stable">
    <img src="https://img.shields.io/badge/go-1.16-blue.svg?style=flat-square&logo=go">
  </a>
  <a href="https://www.python.org/">
    <img src="https://img.shields.io/badge/python-3.6+-blue.svg?style=flat-square&logo=go">
  </a>
   <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">
    <img src="https://img.shields.io/badge/license-GNU-green.svg?style=square&logo=gnu">
   <a href="https://twitter.com/JosueEncinar">
    <img src="https://img.shields.io/badge/author-@JosueEncinar-orange.svg?style=square&logo=twitter">&nbsp;
    <a href="https://twitter.com/topscoder">
    <img src="https://img.shields.io/badge/author-@topscoder-orange.svg?style=square&logo=twitter">
  </a>
</p>


<p align="center">
This script tries to get related domains and/or subdomains by looking at Google Analytics IDs from a URL. First search for ID of Google Analytics in the webpage and then request to <b>builtwith</b> and <b>hackertarget</b> with the ID.</p>

<p align="center">
<b>Note: This is a fork of the original project at github.com/Josue87/AnalyticsRelationships</b>
</p>

<hr/>

 **Note**: It does not work with all websites. It is searched by the following expressions:

* `"www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+"`
* `GTM-[A-Z0-9]+`
* `"UA-\d+-\d+"`

## Installation

Install Golang, then run:

`go install -v github.com/topscoder/analyticsrelationships@latest`

## Usage

This tool can be used in different ways:

1. Pass a single URL using the `-url` flag:
```
analyticsrelationships -url https://www.example.com
```

2. You can also pass URL's as input via STDIN

```
cat urls.txt | analyticsrelationships
```

3. Or a single URL via STDIN

```
echo https://www.example.com | analyticsrelationships
```

## Options

- `-url`: URL of the website to scan the Analytics code from
- `-silent`: Don't print shizzle. Only what matters.

## Contributing

Contributions are welcome! If you find a bug or want to suggest a new feature, please open an issue or submit a pull request.

## License

AnalyticsRelationships is released under the GNU license.