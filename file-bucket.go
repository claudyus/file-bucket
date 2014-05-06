package main

import (
    "fmt"
    "net/http"
    "log"
    "io"
    "os"
    "encoding/json"
    "regexp"
    "strings"
    "path/filepath"
)

var conf Configuration

// test with curl -X POST localhost:1234/upload -F  file=@<FILE>
func BucketRepoHandler(w http.ResponseWriter, req *http.Request) {

    // token validation
    token := strings.Split(req.URL.Path, "/")[1]
    if !conf.bucketExists(token) {
        http.Error(w, "token doesn't exist", 403)
        return
    }

    file_r, handler, err := req.FormFile("file")
    if err != nil {
        http.Error(w, "missed 'file' field in form", 412)
        return
    }

    /* call the pre-push hook (if exists):
     *   if it returns 0 continue
     *   if it returns 1 deny the upload
     *   if it returns 2 override is allowed,
     */
    override := false
    cmd := exec.Command("/etc/file-bucket/hooks/pre-push.sh")
    exit_code := cmd.Run();
    if exit_code.Sys() == 2 {
        override = true
    }
    if exit_code.Sys() == 1 {
        http.Error(w, "upload aborted due to pre-push hook", 409)
        return
    }

    path := filepath.Join(conf.Home, token)
    dst_file := filepath.Join(path, handler.Filename)
    if _, err = os.Stat(dst_file); err == nil {
        http.Error(w, "file exists", 405)
        return
    }

    if err = os.MkdirAll(path, 0700); err != nil {
        http.Error(w, "cannot create bucket dir or writable", 401)
        return
    }

    file_w, err := os.Create(dst_file)
    if err != nil { panic(err) }
    defer file_w.Close()

    fmt.Printf(" * Recieving file %s\n", dst_file)
    if _, err = io.Copy(file_w, file_r); err != nil {
        panic(err)
    }

}

type Configuration struct {
    Buckets []string
    Host string
    Port int
    Home string
}
func (c *Configuration) listenAddr() string {
    if c.Host == "" {
        c.Host = "0.0.0.0"
    }
    if c.Port == 0 {
        c.Port = 1234
    }
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
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
    file, err := os.Open("/etc/file-bucket/config.json")
    if err != nil {
        file, err = os.Open("config.json")
        if err != nil {
            fmt.Println("ERROR cannot read config file")
            return
        }
    }
    decoder := json.NewDecoder(file)
    conf = Configuration{}
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
