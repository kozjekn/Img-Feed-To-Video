# Image feed to video
This project transforms files inside a folder into a video. Used for compatibility with old camera which doesn't support `rtsp` protocol.

## Goals
- [x] Create video from files
- [x] Edit each image so that it has mod date in top left corner
- [x] Add support for arguments
- [x] Remove files once video is created
- [x] Save video to custom location
- [ ] GitHub Release via actions

## Build & run project

Run tidy command (required only once):

```sh
go mod tidy
```

To run project simply run (or better yet use VS code debugger):

```sh
go run . false ./examples ./output/test.avi
```

## Resources
- https://github.com/icza/mjpeg