//go:build enterprise
// +build enterprise

package tbase_audit_log

import (
	"log"
	"path/filepath"
	"regexp"

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

func (w *Watcher) WatchFileCreated(path string) {
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
				regex := regexp.MustCompile(`^postgresql-.*\.csv$`)
				fileName := filepath.Base(event.Name)
				if regex.MatchString(fileName) {
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
