osdtools
========

These are some tools that provide OSD functionality for volume and brightness.

Here is what they look like:

![Volume](screenshot/volume.jpg)
![Mute](screenshot/mute.jpg)
![Brightness](screenshot/brightness.jpg)

They are written in Go, using GTK, and are built with Wayland in mind.
Compiling them takes a while, and the binaries are big. You have been warned.
Install gtk3, gtk-layer-shell and the Go compiler and then you can install
these tools with:

```
go get -u github.com/osdtools/cmd/...
```
