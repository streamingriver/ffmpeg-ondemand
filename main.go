package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/9glt/go-signals"
)

var (
	ffmpeg = &FFmpeg{}

	flagURL  = flag.String("url", "", "URL to stream")
	flagPath = flag.String("path", "", "Path to ts files")
	flagBind = flag.String("bind", ":9999", "Bind address")
)

func resetFFMPEG() {
	input := `-xerror -nostats -nostdin -i {url} -codec copy -map 0:0 -map 0:1 -map_metadata 0 -hls_flags delete_segments -hls_time 10 -segment_list_size 6 -hls_segment_filename file%07d.ts stream.m3u8`
	input = strings.ReplaceAll(input, "{url}", *flagURL)
	argz := strings.Split(input, " ")
	ffmpeg = New(argz, *flagPath)
}

func main() {
	flag.Parse()

	if *flagURL == "" {
		log.Fatal("URL is required")
	}
	if *flagPath == "" {
		log.Fatal("Path is required")
	}
	resetFFMPEG()

	os.MkdirAll(*flagPath, 0755)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if ffmpeg.HitExpired() {
				ffmpeg.Stop()
				ffmpeg.Wait()
				resetFFMPEG()
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Request: %s", r.URL.RequestURI())
		if strings.Contains(r.URL.Path, "stream.m3u8") {
			if ffmpeg.Start() != nil {
				time.Sleep(3 * time.Second)
			}
		}
		ffmpeg.Hit()
		w.Header().Add("Expires", "0")
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	})

	signals.INT(func() {
		if ffmpeg.running {
			ffmpeg.Stop()
		}
		os.Exit(1)
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
	runtime.Goexit()
}
