package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "flag"
    "net/http"
    "os"
    "os/exec"
    "os/signal"
    "path/filepath"
    "regexp"
    "strings"
    "syscall"
)

var conf Configuration


func BucketRepoHandler(w http.ResponseWriter, req *http.Request) {

    /* token validation */
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
     *   if it returns 2 overwrite is allowed,
     */
    overwrite := false
    cmd := exec.Command("/etc/file-bucket/pre-push.sh", token,
        handler.Filename, "0", req.RemoteAddr)  // FIXME find the file size
    err = cmd.Run()
    // command return a I/O error o exitcode != 0
    if err != nil {
        ec, ok := err.(*exec.ExitError);
        if ok {
            if ec.Sys().(syscall.WaitStatus).ExitStatus() == 1 {
                http.Error(w, "upload aborted due to pre-push hook", 409)
                return
            }
            if ec.Sys().(syscall.WaitStatus).ExitStatus() == 2 {
                overwrite = true
            }
        }
    }

    /* ensure that the file doesn't yet exists */
    path := filepath.Join(conf.Home, token)
    dst_file := filepath.Join(path, handler.Filename)
    if _, err = os.Stat(dst_file); err == nil && overwrite == false {
        http.Error(w, "file exists", 405)
        return
    }

    /* create the bucket path if needed or exit */
    if err = os.MkdirAll(path, 0700); err != nil {
        http.Error(w, "cannot create bucket dir or writable", 401)
        return
    }

    /* create, recieve and close the file */
    file_w, err := os.Create(dst_file)
    if err != nil {
        panic(err)
    }

    fmt.Printf(" * Recieving file %s\n", dst_file)
    if _, err = io.Copy(file_w, file_r); err != nil {
        panic(err)
    }
    file_w.Close()

    /* execute the post-push hook and return stdout to client */
    out, err := exec.Command("/etc/file-bucket/post-push.sh", token, dst_file, req.RemoteAddr).Output()
    if err == nil {
        w.Write(out)
    }
}


type Configuration struct {
    Buckets []string
    Host string
    Port int
    Home string
    ActualConfigFile string
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


func signalHandler () {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGUSR2)
    for ; true; {
        sig := <- ch
        log.Printf("Signal %v recieved, re-read bucket lists", sig)
        file, err := os.Open(conf.ActualConfigFile)
        if err == nil {
            decoder := json.NewDecoder(file)
            err = decoder.Decode(&conf)
        }
    }
}


func main() {

    /* parse the command line */
    var config = flag.String("config", "/etc/file-bucket/config.json", "set configuration file")
    flag.Parse()

    /* open, read config file and init the struct */
    conf = Configuration{}
    conf.ActualConfigFile = *config;
    file, err := os.Open(conf.ActualConfigFile)
    if err != nil {
        fmt.Println("ERROR: cannot read config file")
        os.Exit(1)
    }
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&conf)
    if err != nil {
      fmt.Println("error:", err)
    }

    /* setup signal handler */
    go signalHandler();

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
