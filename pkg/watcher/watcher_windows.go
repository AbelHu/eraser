//go:build windows
// +build windows
package watcher

import (
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type windowsFileWatcher struct {
	watcher *fsnotify.Watcher
	path    string
	events  chan FileEvent
	errors  chan error
	done    chan struct{}
}

func newPlatformWatcher(path string) (FileWatcher, error) {
	// On Windows, we'll watch the directory containing the file
	dirPath := filepath.Dir(path)
	filename := filepath.Base(path)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &windowsFileWatcher{
		watcher: fsWatcher,
		path:    path,
		events:  make(chan FileEvent),
		errors:  make(chan error),
		done:    make(chan struct{}),
	}

	if err := fsWatcher.Add(dirPath); err != nil {
		fsWatcher.Close()
		return nil, err
	}

	go w.watch(filename)
	return w, nil
}

func (w *windowsFileWatcher) watch(filename string) {
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
			if filepath.Base(event.Name) == filename {
				w.events <- FileEvent{
					Name: w.path,
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

func (w *windowsFileWatcher) Events() <-chan FileEvent {
	return w.events
}

func (w *windowsFileWatcher) Errors() <-chan error {
	return w.errors
}

func (w *windowsFileWatcher) Close() error {
	close(w.done)
	return w.watcher.Close()
}

func (w *windowsFileWatcher) Reset() error {
	// Windows file watching can be a bit tricky, so we'll implement a simple
	// reset by closing and recreating the watcher
	w.watcher.Close()

	// Wait a moment to let system release resources
	time.Sleep(100 * time.Millisecond)

	// Create new watcher
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	w.watcher = fsWatcher
	dirPath := filepath.Dir(w.path)

	// Add the directory to watch
	return w.watcher.Add(dirPath)
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
