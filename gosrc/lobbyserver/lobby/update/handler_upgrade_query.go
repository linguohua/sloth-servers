package update

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lobbyserver/config"
	"net/http"
	"os"
	// "fmt"
	"github.com/Masterminds/semver"
)

var (
	queryErrorParamQModIsNull = 1
	queryErrorParamModVIsNull = 2
	queryErrorModuleNotExist  = 3
	queryErrorUnmarshalCfg    = 4
)

// Dep 依赖
type Dep struct {
}

// Bundle AssetBundle
type Bundle struct {
	Name string `json:"name"`
	MD5  string `json:"md5"`
	Size int64  `json:"size"`
	Deps []Dep  `json:"deps"`
}

// UpgradeQueryReply 更新查询回复
type UpgradeQueryReply struct {
	Code    int      `json:"code"`
	ABValid bool     `json:"abValid"`
	Name    string   `json:"name"`
	Version string   `json:"version"`
	ABList  []Bundle `json:"abList"`
}

func replyUpgradeQuery(w http.ResponseWriter, reply *UpgradeQueryReply) {
	buf, err := json.Marshal(reply)
	if err != nil {
		log.Println("replyWxLogin, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func isFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func findMatchMaxVersionString(currentVersionStr string, versions []string) string {
	maxVersion := semver.MustParse(currentVersionStr)
	incMajorVersion := maxVersion.IncMajor()
	// log.Printf("maxVersion:%v, incMajorVersion:%v", maxVersion, incMajorVersion)
	for _, ver := range versions {
		version := semver.MustParse(ver)
		// 最大版本相同，才可以通过热更更新
		if version.LessThan(&incMajorVersion) && version.GreaterThan(maxVersion) {
			maxVersion = version
		}
	}

	return maxVersion.Original()
}

func handlerUpgradeQuery(w http.ResponseWriter, r *http.Request) {
	qMod := r.URL.Query().Get("qMod")
	modV := r.URL.Query().Get("modV")

	reply := &UpgradeQueryReply{}

	if qMod == "" {
		reply.Code = queryErrorParamQModIsNull
	}

	if modV == "" {
		reply.Code = queryErrorParamModVIsNull
	}

	if reply.Code != 0 {
		replyUpgradeQuery(w, reply)
		return
	}

	dirPath := config.FileServerPath + "/" + qMod

	isModuleExist, _ := isFileExist(dirPath)
	if !isModuleExist {
		reply.Code = queryErrorModuleNotExist
		replyUpgradeQuery(w, reply)
		return
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	versions := make([]string, 0)
	for _, f := range files {
		versionString := f.Name()
		versions = append(versions, versionString)
	}

	maxVersion := findMatchMaxVersionString(modV, versions)
	if maxVersion == modV {
		reply.Code = 0
		reply.ABValid = false
		replyUpgradeQuery(w, reply)
		return
	}

	log.Println("maxVersion:", maxVersion)
	cfgFilePath := dirPath + "/" + maxVersion + "/cfg.json"
	isCfgFileExist, _ := isFileExist(cfgFilePath)
	if !isCfgFileExist {
		reply.Code = 0
		reply.ABValid = false
		replyUpgradeQuery(w, reply)
		return
	}

	b, err := ioutil.ReadFile(cfgFilePath) // just pass the file name
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(b, reply)
	if err != nil {
		reply.Code = queryErrorUnmarshalCfg
		replyUpgradeQuery(w, reply)
		return
	}

	reply.Code = 0
	reply.ABValid = true

	replyUpgradeQuery(w, reply)
}
