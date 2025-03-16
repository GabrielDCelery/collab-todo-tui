package commands

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type Note struct {
	Path    string
	Name    string
	ModTime int64
}

type NotesRegistered map[string]Note

type NoteCreated struct {
	Path string
	Note Note
}

type NoteRemoved struct {
	Path string
}

func RegisterNotes(dir string, allowedExtensions []string) tea.Cmd {
	return func() tea.Msg {
		notes := map[string]Note{}

		err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			info, err := entry.Info()
			if !slices.Contains(allowedExtensions, filepath.Ext(info.Name())) {
				return nil
			}
			note := Note{
				Path:    path,
				Name:    info.Name(),
				ModTime: info.ModTime().UnixMilli(),
			}
			notes[path] = note
			return nil
		})

		if err != nil {
			return err
		}

		return NotesRegistered(notes)
	}
}

func WaitForNoteCreated(noteCreatedChan chan NoteCreated) tea.Cmd {
	return func() tea.Msg {
		return NoteCreated(<-noteCreatedChan)
	}
}

func WaitForNoteRemoved(noteRemovedChan chan NoteRemoved) tea.Cmd {
	return func() tea.Msg {
		return NoteRemoved(<-noteRemovedChan)
	}
}

func ListenForNoteChanges(ctx context.Context, dir string, allowedExtensions []string, noteCreatedChan chan NoteCreated, noteRemovedChan chan NoteRemoved) tea.Cmd {
	return func() tea.Msg {
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

		watcherErr = watcher.Add(dir)

		if watcherErr != nil {
			return fmt.Errorf("failed to add directory %s to watcher: %w", dir, watcherErr)
		}

		defer watcher.Close()

		for {
			select {
			case <-ctx.Done():
				return nil
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if !slices.Contains(allowedExtensions, filepath.Ext(event.Name)) {
					continue
				}
				switch {
				case event.Has(fsnotify.Create), event.Has(fsnotify.Write):
					note, err := getNoteInfo(event.Name)
					if err != nil {
						return fmt.Errorf("watcher error: %w", err)
					}
					noteCreatedChan <- NoteCreated{
						Path: event.Name,
						Note: note,
					}
				case event.Has(fsnotify.Remove), event.Has(fsnotify.Rename):
					noteRemovedChan <- NoteRemoved{
						Path: event.Name,
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				return fmt.Errorf("watcher error: %w", err)
			}
		}
	}
}

func getNoteInfo(path string) (Note, error) {
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
