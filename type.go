package main

const (
	defaultFileMode = 0644
)

type inputConfig struct {
	Files        []inputFile
	SystemdUnits []inputSystemdUnit
	Users        []inputUser
}

type inputFile struct {
	Path     string
	Contents string
}

type inputSystemdUnit struct {
	Name     string
	Contents string
}

type inputUser struct {
	Name           string
	AuthorizedKeys []string
	PasswordHash   string
}

type ctConfig struct {
	Storage struct {
		Files []ctFile `yaml:"files"`
	} `yaml:"storage"`
	Systemd struct {
		Units []ctSystemdUnit `yaml:"units"`
	} `yaml:"systemd"`
	Passwd struct {
		Users []ctUser `yaml:"users"`
	} `yaml:"passwd"`
}

type ctFile struct {
	Filesystem string `yaml:"filesystem"`
	Path       string `yaml:"path"`
	Mode       *int   `yaml:"mode"`
	Contents   struct {
		Inline string `yaml:"inline"`
	} `yaml:"contents"`
}

type ctSystemdUnit struct {
	Name     string `yaml:"name"`
	Enabled  *bool  `yaml:"enabled"`
	Contents string `yaml:"contents"`
}

type ctUser struct {
	Name              string   `yaml:"name"`
	PasswordHash      *string  `yaml:"password_hash"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}
