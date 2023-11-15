# Image feed to video
This project transforms files inside a folder into a video. Used for compatibility with old camera which doesn't support `rtsp` protocol.

## Goals
- [ ] Create video from files
- [ ] Remove files once video is created
- [ ] Save video to custom location

## Build & run project

Run tidy command (required only once):

```sh
go mod tidy
```

To run project simply run (or better yet use VS code debugger):

```sh
go run .
```

## Resources
- https://github.com/icza/mjpeg