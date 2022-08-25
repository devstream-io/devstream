package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/siddontang/go/log"
)

var Configs []RepositoryConfig

func main() {
	loadConfig()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("webhook", webhook)

	}
	r.Run(":8000")
}


func webhook(c *gin.Context) {

	message,err := c.GetRawData()
	if err != nil{
		log.Info(err.Error())
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(message))
	var form GiteaForm
	if c.ShouldBind(&form) != nil {

		c.JSON(400, gin.H{"message": "invalid parameters"})
		return
	}
	form.Secret=c.GetHeader("secert")
	found := false
	for _, config := range Configs {
		secret:="sha256="+HmacSha256(string(message),config.Secret)
		if secret == form.Secret && config.FullName == form.Repository.FullName && config.Ref == form.Ref {
			cmd := exec.Command("bash", "-c", "cd "+config.Dir+" && git pull")
			if stdout, err := cmd.CombinedOutput(); err != nil {
				c.JSON(500, gin.H{"message": stdout})
				return
			}
			if config.Exec != "" {
				cmd = exec.Command("bash", "-c", config.Exec)
				if stdout, err := cmd.CombinedOutput(); err != nil {
					c.JSON(500, gin.H{"message": stdout})
					return
				}
			}
			found = true
		}
	}
	if !found {
		c.JSON(404, gin.H{"message": "repository full_name not found in config.json"})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
func gitRefelect(git string){
	
}


func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var rawConfigs []RepositoryConfig
	if decoder.Decode(&rawConfigs) != nil {
		log.Info("cannot decode config.json")

		}

	for _, config := range rawConfigs {
		if config.Secret == "" || config.Dir == "" || config.FullName == "" || config.Ref == "" {
			log.Info("wrong format with config.json")
			}
		}

	Configs = rawConfigs

}

func HmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}



type GiteaForm struct {
	Secret   	string 				`json:"secret"`
	Ref        string                `json:"ref" binding:"required"`
	Repository WebhookRepositoryForm `json:"repository" binding:"required"`
}

type WebhookRepositoryForm struct {
	FullName string `json:"full_name" binding:"required"`
}

type RepositoryConfig struct {
	git      string `json:"git" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
	Ref      string `json:"ref" binding:"required"`
	Dir      string `json:"dir" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Exec     string `json:"exec"`
}