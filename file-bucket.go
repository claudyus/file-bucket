package main

import (
    "fmt"
    "net/http"
    "log"
    "io"
    "os"
    "encoding/json"
    "regexp"
)

// test with curl -X POST localhost:1234/upload -F  file=@<FILE>
func BucketRepoHandler(w http.ResponseWriter, req *http.Request) {
    

    file_r, handler, err := req.FormFile("file")
    if err != nil {
        fmt.Println(err)
    }

    file_w, err := os.Create(handler.Filename)
    if err != nil { panic(err) }
    defer file_w.Close()

    _, err = io.Copy(file_w, file_r)
    if err != nil {
        panic(err)
    }

}

type Configuration struct {
    Buckets []string
    host string
    port int
    buckets_home string
}
func (c *Configuration) listenAddr() string {
    if c.host == "" {
        c.host = "0.0.0.0"
    }
    if c.port == 0 {
        c.port = 1234
    }
    return fmt.Sprintf("%s:%d", c.host, c.port)
}
func (c *Configuration) bucketExists(name string) bool {
    for _, element := range c.Buckets {
        if element == name {
            return true
        }
    }
    return false
}

/*
 * http regex handler copied by
 * http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc
 */
type route struct {
    pattern *regexp.Regexp
    handler http.Handler
}

type RegexpHandler struct {
    routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
    h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
    h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range h.routes {
        if route.pattern.MatchString(r.URL.Path) {
            route.handler.ServeHTTP(w, r)
            return
        }
    }
}


func main() {
    /*
     * open read config file and init the struct
     */
    file, err := os.Open("config.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    conf := Configuration{}
    err = decoder.Decode(&conf)
    if err != nil {
      fmt.Println("error:", err)
    }
    fmt.Println(conf)

    /*
     * Setup the webserver goroutine using
     * custom handler with regex support
     */
    handler := &RegexpHandler{}
    handler.HandleFunc(regexp.MustCompile("/.{32}"), BucketRepoHandler)
    err = http.ListenAndServe(conf.listenAddr(), handler)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
