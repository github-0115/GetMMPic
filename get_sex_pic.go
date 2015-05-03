package main
 
import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strings"
    "sync"
    "sync/atomic"
    _ "time"
)
 
const (
    bufferSize = 64 * 1024 //写图片文件的缓冲区大小
)
 
var (
    numPoller     = flag.Int("p", 5, "page loader num")
    numDownloader = flag.Int("d", 10, "image downloader num")
    savePath      = flag.String("s", "./downloads/", "save path")
    imgExp        = regexp.MustCompile(`<a\s+class="img"\s+href="[a-zA-Z0-9_\-/:\.%?=]+">[\r\n\t\s]*<img\s+src="([^"'<>]*)"\s*/?>`)
)
 
type image struct {
    url      string
    filename string
}
 
type sexyContext struct {
    pollerDone   chan struct{}
    images       map[string]int
    imagesLock   *sync.Mutex
    imageChan    chan *image
    pageIndex    int32
    rootURL      string
    done         bool
    imageCounter int32
    okCounter    int32
}
 
func main() {
    flag.Parse()
    ctx := &sexyContext{
        pollerDone: make(chan struct{}),
        images:     make(map[string]int),
        imagesLock: &sync.Mutex{},
        imageChan:  make(chan *image, 100),
        pageIndex:  0,
        rootURL:    "http://me2-sex.lofter.com/tag/美女摄影",
    }
    os.MkdirAll(*savePath, 0777)
    ctx.start()
 
}
 
func (ctx *sexyContext) start() {
 
    waits := sync.WaitGroup{}
    for i := 0; i < *numDownloader; i++ {
        waits.Add(1)
        go func() {
            ctx.downloadImage()
            waits.Done()
        }()
    }
    waits2 := sync.WaitGroup{}
    for i := 0; i < *numPoller; i++ {
        waits2.Add(1)
        go func() {
            ctx.downloadPage()
            waits2.Done()
        }()
    }
 
    waits2.Wait()
    fmt.Println("poller done")
    close(ctx.pollerDone)
 
    waits.Wait()
    fmt.Printf("fetch done get %d ok %d\n", ctx.imageCounter, ctx.okCounter)
}
 
func (ctx *sexyContext) downloadPage() {
 
    for {
        p := atomic.AddInt32(&ctx.pageIndex, 1)
 
        url := fmt.Sprintf("%s?page=%d", ctx.rootURL, p)
        fmt.Printf("download page %s\n", url)
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf("failed to load url %s with error %v", url, err)
        } else {
 
            body, err := ioutil.ReadAll(resp.Body)
            resp.Body.Close()
            if err != nil {
                fmt.Printf("failed to load url %s with error %v", url, err)
            } else {
                if ctx.parsePage(body) {
                    break
                }
            }
        }
    }
}
 
func (ctx *sexyContext) parsePage(body []byte) bool {
    //fmt.Printf("%s\n", string(body))
    idx := imgExp.FindAllSubmatchIndex(body, -1)
    if idx != nil {
 
        for _, n := range idx {
            imgeUrl := strings.TrimSpace(string(body[n[2]:n[3]]))
            filename := url.QueryEscape(imgeUrl)
            //fmt.Printf("%s\n", filename)
            image := &image{url: imgeUrl, filename: filename}
            atomic.AddInt32(&ctx.imageCounter, 1)
            ctx.imageChan <- image
        }
        return false
    }
    return true
}
 
func (ctx *sexyContext) downloadImage() {
    isDone := false
    for !isDone {
        select {
        case <-ctx.pollerDone:
            if len(ctx.imageChan) == 0 {
                isDone = true
                fmt.Println("poller done and quit")
            }
        case image := <-ctx.imageChan:
            fmt.Printf("start download %s\n", image.url)
            atomic.AddInt32(&ctx.okCounter, 1)
            resp, err := http.Get(image.url)
            if err != nil {
                fmt.Printf("failed to load url %s with error %v\n", image.url, err)
            } else {
                go func() {
                    defer resp.Body.Close()
                    saveFile := *savePath + image.filename //path.Base(imgUrl)
 
                    img, err := os.Create(saveFile)
                    if err != nil {
                        fmt.Print(err)
 
                    } else {
                        defer img.Close()
 
                        imgWriter := bufio.NewWriterSize(img, bufferSize)
 
                        _, err = io.Copy(imgWriter, resp.Body)
                        if err != nil {
                            fmt.Print(err)
 
                        }
                        imgWriter.Flush()
                    }
                }()
            }
        }
    }
 
}
