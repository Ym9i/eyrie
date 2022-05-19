package config

import (
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

var iniConf *ini.File

var (
	AliRefreshToken string
	DownloadPath    string
	Formats         []string
	CookieFile      string
	UserAgent       string
	InitialBookID   int
	ProgressFile    string

	Timeout int
	Retry   int
	Thread  int
	Rename  bool
	Debug   bool

	TestShareIdAnnomous string
	TestShareId         string
	TestShareUrl        string
	TestSharePwd        string
)

func init() {
	var err error

	envPath := "env.ini"
	// look for env.ini
	times := 0
	for {
		times++
		_, err = os.Stat(envPath)
		if err != nil {
			if times > 5 {
				panic("env.ini file not found")
			}
			envPath = "../" + envPath
		} else {
			break
		}
	}

	iniConf, err = ini.Load(envPath)

	if err != nil {
		panic(err)
	}

	AliRefreshToken = loadString(iniConf, "EYRIE", "REFRESH_TOKEN", "")
	DownloadPath = loadString(iniConf, "EYRIE", "DOWNLOAD_PATH", "")
	Formats = loadStringL(iniConf, "EYRIE", "FORMATS", []string{"EPUB"})
	CookieFile = loadString(iniConf, "EYRIE", "COOKIE_FILE", "cookie")
	UserAgent = loadString(iniConf, "EYRIE", "USER_AGENT", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")

	InitialBookID = loadInt(iniConf, "EYRIE", "INITIAL_BOOK_ID", 1)

	ProgressFile = loadString(iniConf, "EYRIE", "PROCESS_FILE", "process")

	Timeout = loadInt(iniConf, "EYRIE", "TIMEOUT", 5000)
	Retry = loadInt(iniConf, "EYRIE", "TETRY", 4)
	Thread = loadInt(iniConf, "EYRIE", "THREAD", 3)
	Rename = loadBool(iniConf, "EYRIE", "RENAME", true)
	Debug = loadBool(iniConf, "EYRIE", "DEBUG", true)

	TestShareIdAnnomous = loadString(iniConf, "EYRIE", "TEST_SHARE_ID_ANONYMOUS", "")
	TestShareId = loadString(iniConf, "EYRIE", "TEST_SHARE_ID", "")
	TestShareUrl = loadString(iniConf, "EYRIE", "TEST_SHARE_URL", "")
	TestSharePwd = loadString(iniConf, "EYRIE", "TEST_SHARE_PWD", "")
}

func loadString(iniF *ini.File, section string, key string, defaultValue string) string {
	v := iniF.Section(section).Key(key).String()

	if v == "" {
		v = defaultValue
	}

	return v
}

func loadStringL(iniF *ini.File, section string, key string, defaultValue []string) []string {
	s := loadString(iniF, section, key, "")

	if s == "" {
		return defaultValue
	}

	sL := strings.Split(s, ",")

	return sL
}

func loadInt(iniF *ini.File, section string, key string, defaultValue int) int {
	v := iniF.Section(section).Key(key).MustInt()

	if v == 0 {
		v = defaultValue
	}

	return v
}

func loadBool(iniF *ini.File, section string, key string, defaultValue bool) bool {
	v := iniF.Section(section).Key(key).MustBool()

	return v
}
