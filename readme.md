osdtools
========

These are some tools that provide OSD functionality for volume and brightness.

Here is what they look like:

![Volume](screenshots/volume.jpg)
![Mute](screenshots/mute.jpg)
![Brightness](screenshots/brightness.jpg)

They are written in Go, using GTK, and are built with Wayland in mind.
Compiling them takes a while, and the binaries are big. You have been warned.
Install gtk3, gtk-layer-shell and the Go compiler and then you can install
these tools with:

```
go get -u github.com/osdtools/cmd/...
```
