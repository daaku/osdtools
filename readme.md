osdtools
========

These are some tools that provide OSD functionality for volume and brightness.

Here is what they look like:

![Volume](screenshots/volume.jpg =321x321)
![Mute](screenshots/mute.jpg =304x304)
![Brightness](screenshots/brightness.jpg =319x319)

They are written in Go, using GTK, and are built with Wayland in mind.
Compiling them takes a while, and the binaries are big. You have been warned.
Install gtk3, gtk-layer-shell and the Go compiler and then you can install
these tools with:

```
go get -u github.com/osdtools/cmd/...
```
