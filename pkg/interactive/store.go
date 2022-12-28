package interactive

import (
	"os"
	"path"
	"strings"
)

type Store struct {
	directory    string
	storeHistory bool
	storeVsm     bool
}

func NewStore() *Store {
	return &Store{
		directory:    path.Join(os.Getenv("HOME"), ".bb"),
		storeHistory: true,
		storeVsm:     true,
	}
}

func (s *Store) Init() error {
	dir := s.directory
	if s.storeVsm {
		dir = path.Join(s.directory, "vsm")
	}
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Directory() string {
	return s.directory
}

func (s *Store) SetVsmStorage(vsm bool) {
	s.storeVsm = vsm
}

func (s *Store) SetHistoryStorage(history bool) {
	s.storeVsm = history
}

func (s *Store) StoreHistoryEntry(entry string) error {
	if s.storeHistory {
		file, err := os.OpenFile(path.Join(s.directory, ".bb_history"),
			os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(
			strings.TrimSpace(entry) + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) StoreVsmIndex(index string, bin []byte) error {
	if s.storeVsm {
		file, err := os.OpenFile(path.Join(s.directory, "vsm", index+".vsm"),
			os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(bin)
		if err != nil {
			return err
		}
	}
	return nil
}
