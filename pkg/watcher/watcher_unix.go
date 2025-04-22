//go:build !windows
// +build !windows

package watcher

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

type unixFileWatcher struct {
	watcher *fsnotify.Watcher
	path    string
	events  chan FileEvent
	errors  chan error
	done    chan struct{}
}

func newPlatformWatcher(path string) (FileWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &unixFileWatcher{
		watcher: fsWatcher,
		path:    path,
		events:  make(chan FileEvent),
		errors:  make(chan error),
		done:    make(chan struct{}),
	}

	if err := fsWatcher.Add(path); err != nil {
		fsWatcher.Close()
		return nil, err
	}

	go w.watch()
	return w, nil
}

func (w *unixFileWatcher) watch() {
	defer close(w.events)
	defer close(w.errors)

	for {
		select {
		case <-w.done:
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Name == w.path {
				w.events <- FileEvent{
					Name: event.Name,
					Op:   convertOp(event.Op),
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.errors <- err
		}
	}
}

func (w *unixFileWatcher) Events() <-chan FileEvent {
	return w.events
}

func (w *unixFileWatcher) Errors() <-chan error {
	return w.errors
}

func (w *unixFileWatcher) Close() error {
	close(w.done)
	return w.watcher.Close()
}

func (w *unixFileWatcher) Reset() error {
	// Remove existing watch if it exists
	_ = w.watcher.Remove(w.path)

	// Check if file exists before adding watch
	if _, err := os.Stat(w.path); os.IsNotExist(err) {
		return err
	}

	return w.watcher.Add(w.path)
}

func convertOp(op fsnotify.Op) Op {
	var result Op
	if op&fsnotify.Create != 0 {
		result |= Create
	}
	if op&fsnotify.Write != 0 {
		result |= Write
	}
	if op&fsnotify.Remove != 0 {
		result |= Remove
	}
	if op&fsnotify.Rename != 0 {
		result |= Rename
	}
	if op&fsnotify.Chmod != 0 {
		result |= Chmod
	}
	return result
}
