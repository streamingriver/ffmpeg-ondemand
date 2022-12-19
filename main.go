package main

import (
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
)

func getFullPath() string {
	return os.Getenv("APP_ROOT") + "/" + os.Getenv("APP_NAME")
}

func resetFFMPEG() {
	input := `-xerror -nostats -nostdin -i {url} -codec copy -map 0:0 -map 0:1 -map_metadata 0 -hls_flags delete_segments -hls_time 10 -segment_list_size 6 -hls_segment_filename file%07d.ts stream.m3u8`
	input = strings.ReplaceAll(input, "{url}", os.Args[1])
	argz := strings.Split(input, " ")
	ffmpeg = New(argz, getFullPath())
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("URL is required")
	}
	if len(os.Args) > 1 && os.Args[1] == "" {
		log.Fatal("URL is required")
	}
	if os.Getenv("APP_ROOT") == "" {
		log.Fatal("Output ROOT is required")
	}
	if os.Getenv("APP_NAME") == "" {
		log.Fatal("Name is required")
	}
	os.MkdirAll(getFullPath(), 0777)

	resetFFMPEG()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if ffmpeg.running {
				if ffmpeg.HitExpired() {
					ffmpeg.Stop()
					ffmpeg.Wait()
					resetFFMPEG()
				}
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
		http.FileServer(http.Dir(getFullPath())).ServeHTTP(w, r)
	})

	signals.INT(func() {
		if ffmpeg.running {
			ffmpeg.Stop()
		}
		os.Exit(1)
	})

	log.Fatal(http.ListenAndServe(os.Getenv("APP_BIND"), nil))
	runtime.Goexit()
}
