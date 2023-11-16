# Image feed to video
This project transforms files inside a folder into a video. Used for compatibility with old camera which doesn't support modern(ish) video streaming protocols eg. `rtsp` protocol. and only supports `ftp` file feed.

## Features
- Create video from files
- Edit each image so that it has mod date in top left corner
- Add support for arguments
- Remove files once video is created
- Save video to custom location
- Build via GitHub actions

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