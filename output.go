package main

func intPtr(n int) *int {
	return &n
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func convert(src *inputConfig) *ctConfig {
	var files []ctFile
	for _, srcFile := range src.Files {
		dstFile := convertFile(srcFile)
		files = append(files, dstFile)
	}

	var units []ctSystemdUnit
	for _, srcUnit := range src.SystemdUnits {
		dstUnit := convertSystemdUnit(srcUnit)
		units = append(units, dstUnit)
	}

	var users []ctUser
	for _, srcUser := range src.Users {
		dstUser := convertUser(srcUser)
		users = append(users, dstUser)
	}

	var dst ctConfig
	dst.Storage.Files = files
	dst.Systemd.Units = units
	dst.Passwd.Users = users
	return &dst
}

func convertFile(src inputFile) ctFile {
	dst := ctFile{
		Filesystem: "root",
		Path:       src.Path,
		Mode:       intPtr(defaultFileMode),
	}
	dst.Contents.Inline = src.Contents
	return dst
}

func convertSystemdUnit(src inputSystemdUnit) ctSystemdUnit {
	dst := ctSystemdUnit{
		Name:     src.Name,
		Enabled:  boolPtr(true),
		Contents: src.Contents,
	}
	return dst
}

func convertUser(src inputUser) ctUser {
	dst := ctUser{
		Name:              src.Name,
		SSHAuthorizedKeys: src.AuthorizedKeys,
	}
	if len(src.PasswordHash) != 0 {
		dst.PasswordHash = stringPtr(src.PasswordHash)
	}
	return dst
}
