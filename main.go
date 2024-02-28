package main

import (
	"embed"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed index.html favicon.ico
var staticFiles embed.FS

var (
	jsdelivr = 0
	// whiteList = []string{}
	// blackList = []string{}
	passList = []string{}
	exp1     = regexp.MustCompile(`^(?:https?://)?github\.com/(?P<author>.+?)/(?P<repo>.+?)/(?:releases|archive)/.*$`)
	exp2     = regexp.MustCompile(`^(?:https?://)?github\.com/(?P<author>.+?)/(?P<repo>.+?)/(?:blob|raw)/.*$`)
	exp3     = regexp.MustCompile(`^(?:https?://)?github\.com/(?P<author>.+?)/(?P<repo>.+?)/(?:info|git-).*$`)
	exp4     = regexp.MustCompile(`^(?:https?://)?raw\.(?:githubusercontent|github)\.com/(?P<author>.+?)/(?P<repo>.+?)/.+?/.+$`)
	exp5     = regexp.MustCompile(`^(?:https?://)?gist\.(?:githubusercontent|github)\.com/(?P<author>.+?)/.+?/.+$`)
	exp6     = regexp.MustCompile(`(\.com/.*?/.+?)/(.+?/)`)
)

func main() {
	var host string
	var port string

	flag.StringVar(&host, "host", "0.0.0.0", "Host address")
	flag.StringVar(&port, "port", "80", "Port number")
	flag.Parse()

	hostRegExp := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	if !hostRegExp.MatchString(host) {
		log.Println("Invalid host: " + host)
		os.Exit(1)
	}

	portRegExp := regexp.MustCompile(`^\d+$`)
	if !portRegExp.MatchString(port) {
		log.Println("Invalid port: " + port)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		q := c.Query("q")
		if q != "" {
			c.Redirect(http.StatusFound, "/"+q)
			return
		}
		data, err := staticFiles.ReadFile("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read index HTML from embed.FS: %v", err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		data, err := staticFiles.ReadFile("favicon.ico")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read favicon from embed.FS: %v", err)
			return
		}
		c.Data(http.StatusOK, "image/vnd.microsoft.icon", data)
	})

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path[1:]
		handler(c, path)
	})

	log.Println("Server started at " + host + ":" + port)
	router.Run(host + ":" + port)
}

func handler(c *gin.Context, u string) {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	// u = strings.Replace(u, "s:/", "s://", 1) // Fix for double slash issue

	passBy := false
	match, _ := checkURL(u)
	if match {
		for _, i := range passList {
			if strings.Contains(u, i) {
				passBy = true
				break
			}
		}
	} else {
		c.String(http.StatusForbidden, "Invalid input.")
		return
	}

	if (jsdelivr > 0 || passBy) && exp2.MatchString(u) {
		u = strings.Replace(u, "/blob/", "@", 1)
		u = strings.Replace(u, "github.com", "cdn.jsdelivr.net/gh", 1)
		c.Redirect(http.StatusFound, u)
	} else if (jsdelivr > 0 || passBy) && exp4.MatchString(u) {
		u = exp6.ReplaceAllString(u, "$1@$2")
		u = strings.Replace(u, "raw.githubusercontent.com", "cdn.jsdelivr.net/gh", 1)
		c.Redirect(http.StatusFound, u)
	} else {
		if exp2.MatchString(u) {
			u = strings.Replace(u, "/blob/", "/raw/", 1)
		}
		proxy(c, u)
	}
}

func proxy(c *gin.Context, u string) {
	client := &http.Client{}
	req, err := http.NewRequest(c.Request.Method, u, nil)
	if err != nil {
		log.Println("Failed to create request: ", err)
		c.String(http.StatusInternalServerError, "Server error: %v", err)
		return
	}

	copyHeader(c.Request.Header, &req.Header)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request: ", err)
		c.String(http.StatusInternalServerError, "Server error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Status(resp.StatusCode)
		return
	}

	header := c.Writer.Header()
	copyHeader(resp.Header, &header)
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func copyHeader(src http.Header, dest *http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dest.Add(k, v)
		}
	}
}

func checkURL(u string) (bool, []string) {
	var allExps = []*regexp.Regexp{
		exp1, exp2, exp3, exp4, exp5,
	}

	for _, exp := range allExps {
		if exp.MatchString(u) {
			return true, exp.FindStringSubmatch(u)
		}
	}
	return false, nil
}
