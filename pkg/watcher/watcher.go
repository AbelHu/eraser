package watcher

// FileEvent represents a file system event
type FileEvent struct {
	Name string
	Op   Op
}

// Op represents the type of file system operation
type Op uint32

const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// IsDelete checks if the event is a delete event
func (e FileEvent) IsDelete() bool {
	return e.Op&Remove != 0
}

// IsRename checks if the event is a rename event
func (e FileEvent) IsRename() bool {
	return e.Op&Rename != 0
}

// FileWatcher provides a platform-independent interface for watching files
type FileWatcher interface {
	// Events returns a channel for receiving file events
	Events() <-chan FileEvent
	// Errors returns a channel for receiving errors
	Errors() <-chan error
	// Close stops watching and closes all channels
	Close() error
	// Reset reestablishes the watch on the file
	Reset() error
}

// NewFileWatcher creates a platform-specific file watcher implementation
// for the provided file path
func NewFileWatcher(path string) (FileWatcher, error) {
	return newPlatformWatcher(path)
}
