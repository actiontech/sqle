//go:build enterprise
// +build enterprise

package tbase_audit_log

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

type Watcher struct {
	l        *logrus.Entry
	NewFiles chan string
}

func NewWatcher() *Watcher {
	return &Watcher{
		l:        logrus.WithField("scanner", "tbase-audit-log"),
		NewFiles: make(chan string, 1),
	}
}

func (w *Watcher) WatchFileCreated(path, fileNameFormat string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				fileName := filepath.Base(event.Name)
				isMatch, err := filepath.Match(fileNameFormat, fileName)
				if err != nil {
					w.l.Error(err)
				}
				if isMatch {
					w.NewFiles <- event.Name
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			w.l.Errorf("watch file create failed, error: %v", err)
		}
	}
}
