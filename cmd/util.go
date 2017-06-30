package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NormalizePath(cp, fp string) string {
	dir := filepath.Dir(cp)

	if fp != "" && !filepath.IsAbs(fp) {
		fp = filepath.Join(dir, fp)
	}

	return fp
}

func ReadSecret(key string, env bool) ([]byte, error) {
	var (
		buf []byte
		err error
	)

	if env {
		ev := os.Getenv(key)
		if ev == "" {
			return nil, fmt.Errorf("environment variable does not exist: %s", key)
		}

		buf = []byte(ev)
	} else {
		buf, err = ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func WriteSecret(key string, data []byte, env bool) (int, error) {
	if env {
		err := os.Setenv(key, string(data))
		if err != nil {
			return 0, err
		}
	} else {
		dir := filepath.Dir(key)

		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				return 0, fmt.Errorf("unable to create directory: %s\n", dir)
			}
		}

		err = ioutil.WriteFile(key, data, 0600)
		if err != nil {
			return 0, err
		}
	}

	return len(data), nil
}
