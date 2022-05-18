package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Secret []byte

func InitSecret() {
	conf, _ := os.Open("./confs/secret.json")
	defer conf.Close()
	value, _ := ioutil.ReadAll(conf)
	json.Unmarshal([]byte(value), &map[string][]byte{"secret": Secret})
}
