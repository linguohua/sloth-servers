package update
import (
	"net/http"
	"encoding/json"
	"lobbyserver/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"fmt"
)

var (
	queryErrorParamQModIsNull = 1
	queryErrorParamModVIsNull = 2
	queryErrorModuleNotExist = 3

)

// UpgradeQueryReply 更新查询回复
type UpgradeQueryReply struct {
	Code int `json:"code"`
	ABValid bool `json:"abValid"`
}

func replyUpgradeQuery(w http.ResponseWriter, reply UpgradeQueryReply) {
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
	maxVersion := parseVersionString(currentVersionStr)
	for _, ver := range versions {
		version := parseVersionString(ver)
		// 最大版本相同，才可以通过热更更新
		if version.bigVer == maxVersion.bigVer {
			if (version.middleVer > maxVersion.middleVer) {
				maxVersion = version
			} else {
				if version.middleVer == maxVersion.middleVer && version.smallVer > maxVersion.smallVer{
					maxVersion = version
				}
			}
		}
	}

	return fmt.Sprintf("%d.%d.%d", maxVersion.bigVer, maxVersion.middleVer, maxVersion.smallVer)
}

	//
	// csVer := r.URL.Query().Get("csVer")
	// lobbyVer := r.URL.Query().Get("lobbyVer")
	// operatingSystem := r.URL.Query().Get("operatingSystem")
	// operatingSystemFamily := r.URL.Query().Get("operatingSystemFamily")
	// deviceUniqueIdentifier := r.URL.Query().Get("deviceUniqueIdentifier")
	// deviceName := r.URL.Query().Get("deviceName")
	// deviceModel := r.URL.Query().Get("deviceModel")
	// network := r.URL.Query().Get("network")

func handlerUpgradeQuery(w http.ResponseWriter, r *http.Request) {
	qMod := r.URL.Query().Get("qMod")
	modV := r.URL.Query().Get("modV")

	reply := UpgradeQueryReply{}

	if (qMod == "") {
		reply.Code = queryErrorParamQModIsNull
	}

	if (modV == "") {
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

	cfgFilePath := dirPath + "/" + maxVersion + "/cfg.json"
	isCfgFileExist, _ := isFileExist(cfgFilePath)
	if !isCfgFileExist {
		reply.Code = 0
		reply.ABValid = false
		replyUpgradeQuery(w, reply)
	}

 	b, err := ioutil.ReadFile(cfgFilePath) // just pass the file name
    if err != nil {
        log.Println(err)
    }

    w.Write(b)
}