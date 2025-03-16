package notes

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/fsnotify/fsnotify"
)

type Note struct {
	Path    string
	Name    string
	ModTime int64
}

type NotesWatcher struct {
	dir               string
	allowedExtensions []string
	Notes             map[string]Note
}

func NewNotesWatcher(dir string, allowedExtensions []string) *NotesWatcher {
	return &NotesWatcher{
		dir:               dir,
		allowedExtensions: allowedExtensions,
		Notes:             map[string]Note{},
	}
}

func (n *NotesWatcher) Start(ctx context.Context) error {
	if err := n.getAllNotes(); err != nil {
		return err
	}
	if err := n.watchDir(ctx); err != nil {
		return err
	}
	return nil
}

func (n *NotesWatcher) getNoteInfo(path string) (Note, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Note{}, fmt.Errorf("failed to get file info: %w", err)
	}
	note := Note{
		Path:    path,
		Name:    info.Name(),
		ModTime: info.ModTime().UnixMilli(),
	}
	return note, nil
}

func (n *NotesWatcher) getAllNotes() error {
	err := filepath.WalkDir(n.dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if !slices.Contains(n.allowedExtensions, filepath.Ext(info.Name())) {
			return nil
		}
		note := Note{
			Path:    path,
			Name:    info.Name(),
			ModTime: info.ModTime().UnixMilli(),
		}
		n.Notes[path] = note
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (n *NotesWatcher) watchDir(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	var watcherErr error
	defer func() {
		if watcherErr != nil {
			watcher.Close()
		}
	}()

	watcherErr = watcher.Add(n.dir)
	if watcherErr != nil {
		return fmt.Errorf("failed to add directory %s to watcher: %w", n.dir, watcherErr)
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !slices.Contains(n.allowedExtensions, filepath.Ext(event.Name)) {
					continue

				}
				switch {
				case event.Has(fsnotify.Create), event.Has(fsnotify.Write):
					note, err := n.getNoteInfo(event.Name)
					if err != nil {
						log.Printf("watcher error: %v\n", err)
					}
					n.Notes[event.Name] = note
				case event.Has(fsnotify.Remove), event.Has(fsnotify.Rename):
					delete(n.Notes, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("watcher error: %v\n", err)
			}
		}
	}()

	return nil
}
