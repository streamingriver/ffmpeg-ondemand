# FFMPEG on Demand

usage: 
```
docker run --rm -it --name test-container -e APP_NAME="channel-name" -p 9999:9999 ghcr.io/streamingriver/ffmpeg-ondemand:mai http://url/to/channel/main.m3u8
```