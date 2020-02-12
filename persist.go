package main

import (
	"os"
	"path/filepath"
	"runtime"

	"gitlab.com/NebulousLabs/Sia/persist"
)

var (
	persistMetadata = persist.Metadata{
		Header:  "SkySync Persistence",
		Version: "v0.1.0",
	}

	persistFileName = filepath.Join(skySyncPersistDir(), "skysync.json")
)

type (
	persistedFile struct {
		Filename string `json:"filename"`
		SkyLink  string `json:"skylink"`
	}
	persistence struct {
		Files []persistedFile `json:"files"`
	}
)

// skySyncPersistDir returns the directory that they skysync persistence will be
// saved
func skySyncPersistDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "SkySync")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "SkySync")
	default:
		return filepath.Join(os.Getenv("HOME"), ".skysync")
	}
}

// load loads the SkySync's persistence from disk
func (ss *SkySync) load() error {
	var data persistence
	err := persist.LoadJSON(persistMetadata, &data, persistFileName)
	if os.IsNotExist(err) {
		err := os.MkdirAll(skySyncPersistDir(), 0700)
		if err != nil {
			return err
		}
		err = ss.save()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	for _, file := range data.Files {
		ss.skyfiles[file.Filename] = file.SkyLink
	}
	return nil
}

func (ss *SkySync) persistData() persistence {
	var data persistence
	for file, skylink := range ss.skyfiles {
		data.Files = append(data.Files, persistedFile{
			Filename: file,
			SkyLink:  skylink,
		})
	}
	return data
}

// save saves the SkySync's persistence to disk
func (ss *SkySync) save() error {
	return persist.SaveJSON(persistMetadata, ss.persistData(), persistFileName)
}
