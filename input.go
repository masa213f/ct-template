package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func readTemplateDir(tmplDir string) (*inputConfig, error) {
	var files []inputFile
	var units []inputSystemdUnit
	var users []inputUser

	subDirs, err := ioutil.ReadDir(tmplDir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, dir := range subDirs {
		if !dir.IsDir() {
			continue
		}

		path := filepath.Join(tmplDir, dir.Name())
		switch dir.Name() {
		case "file":
			f, err := readFileDir(path)
			if err != nil {
				return nil, err
			}
			files = f

		case "systemd":
			u, err := readSystemdDir(path)
			if err != nil {
				return nil, err
			}
			units = u

		case "user":
			u, err := readUserDir(path)
			if err != nil {
				return nil, err
			}
			users = u
		}
	}

	input := inputConfig{
		Files:        files,
		SystemdUnits: units,
		Users:        users,
	}
	return &input, nil
}

func readFileDir(root string) ([]inputFile, error) {
	var files []inputFile
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		ent := inputFile{
			Path:     strings.Replace(path, root, "", 1),
			Contents: string(buf),
		}
		files = append(files, ent)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}

// ref: https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/chap-Managing_Services_with_systemd#tabl-Managing_Services_with_systemd-Introduction-Units-Types
var unitTypes = []string{
	".service",
	".target",
	".automount",
	".device",
	".mount",
	".path",
	".scope",
	".slice",
	".snapshot",
	".socket",
	".swap",
	".timer",
}

func validateFileExtension(name string) bool {
	ext := filepath.Ext(name)
	for _, e := range unitTypes {
		if e == ext {
			return true
		}
	}
	return false
}

func readSystemdDir(root string) ([]inputSystemdUnit, error) {
	var units []inputSystemdUnit
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !validateFileExtension(info.Name()) {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		ent := inputSystemdUnit{
			Name:     info.Name(),
			Contents: string(buf),
		}
		units = append(units, ent)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return units, err
}

func readAuthorizedKeys(path string) ([]string, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := strings.TrimSpace(scanner.Text())
		if len(key) == 0 {
			continue
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func readPasswordHash(path string) (string, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hash := strings.TrimSpace(scanner.Text())
		return hash, nil
	}
	// Should an error be returned?
	return "", nil
}

func readUserDir(root string) ([]inputUser, error) {
	var users []inputUser

	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		keysFile := filepath.Join(root, dir.Name(), "authorized_keys")
		keys, err := readAuthorizedKeys(keysFile)
		if err != nil {
			return nil, err
		}

		hashFile := filepath.Join(root, dir.Name(), "password_hash")
		hash, err := readPasswordHash(hashFile)
		if err != nil {
			return nil, err
		}

		ent := inputUser{
			Name:           dir.Name(),
			AuthorizedKeys: keys,
			PasswordHash:   hash,
		}
		users = append(users, ent)
	}
	return users, err
}
