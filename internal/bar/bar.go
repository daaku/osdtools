package bar

import (
	"net"
	"net/rpc"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/daaku/osdtools/internal/imagewindow"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

const styles = `
window {
	background: rgba(5%, 5%, 5%, 0.75);
}
.icon {
	font-size: 128px;
	padding: 16px;
}
.progress trough,
.progress trough progress {
	min-height: 20px;
	border-radius: 0;
}
.progress trough progress {
	min-height: 20px;
}
`

func socketPath(name string) string {
	u, _ := user.Current()
	home := ""
	if u == nil {
		home = os.Getenv("HOME")
	} else {

		home = u.HomeDir
	}
	return filepath.Join(home, name)
}

// ImageBarServer to handle RPC calls.
type ImageBarServer struct {
	pct chan int
}

// Set sets the pct.
func (s *ImageBarServer) Set(v int, ok *struct{}) error {
	s.pct <- v
	return nil
}

func sendPct(socketPath string, pct int) error {
	client, err := rpc.Dial("unix", socketPath)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := client.Call("ImageBarServer.Set", pct, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func startServer(socketPath string) (chan int, net.Listener, error) {
	pct := make(chan int)
	if err := rpc.Register(&ImageBarServer{pct: pct}); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	go rpc.Accept(lis)
	return pct, lis, nil
}

// Run a bar application.
func Run(pct int, icon func(int) string, socketName string) error {
	const timeout = 5 * time.Second
	socketPath := socketPath(socketName)

	if err := sendPct(socketPath, pct); err == nil {
		return nil
	}

	pctCh, lis, err := startServer(socketPath)
	if err != nil {
		return err
	}
	defer lis.Close()

	gtk.Init(nil)

	if err := imagewindow.ConfigureStyles(styles); err != nil {
		return err
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return errors.WithStack(err)
	}
	win.SetTitle("OSD")
	imagewindow.LayerInit(win)
	imagewindow.LayerSetOverlay(win)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return errors.WithStack(err)
	}

	iconLabel, err := gtk.LabelNew(icon(pct))
	if err != nil {
		return errors.WithStack(err)
	}
	if err := imagewindow.AddClass(iconLabel, "icon"); err != nil {
		return err
	}
	vbox.Add(iconLabel)

	progress, err := gtk.ProgressBarNew()
	if err != nil {
		return errors.WithStack(err)
	}
	progress.SetFraction(float64(pct) / 100)
	if err := imagewindow.AddClass(progress, "progress"); err != nil {
		return err
	}
	vbox.Add(progress)

	win.Add(vbox)
	win.ShowAll()

	go func() {
		quitter := time.NewTimer(timeout)
		for {
			select {
			case <-quitter.C:
				_, _ = glib.IdleAdd(gtk.MainQuit)
				break
			case pct := <-pctCh:
				quitter.Reset(timeout)
				iconLabel.SetText(icon(pct))
				progress.SetFraction(float64(pct) / 100)
			}
		}
	}()

	gtk.Main()
	return nil
}
