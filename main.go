package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

var c Config

type Config struct {
	AllowedHosts []string `json:"allowed_hosts"`
	ApiKey       string   `json:"apikey"`
	Port         int      `json:"port"`
}

func generateApiKey(length int) string {
	var output []byte = make([]byte, length)
	var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < length; i++ {
		output[i] = chars[rand.Intn(len(chars))]
	}
	return string(output)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if token == c.ApiKey {
		if len(c.AllowedHosts) > 0 {
			allowed := false
			host := strings.Split(r.RemoteAddr, ":")[0]
			for _, s := range c.AllowedHosts {
				if host == s {
					allowed = true
				}
			}
			if !allowed {
				log.Printf("Host %s not in allowed list %v", host, c.AllowedHosts)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("shutdown", "/s", "/t", "0")
		case "linux":
			cmd = exec.Command("shutdown", "-h", "now")
		default:
			fmt.Printf("OS %s not supported", runtime.GOOS)
		}

		b, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

func ReadConfig(location string) (Config, error) {

	b, err := os.ReadFile(location)
	if err != nil {
		c = Config{AllowedHosts: []string{}, ApiKey: generateApiKey(32), Port: 7000}
		b, _ = json.MarshalIndent(c, "", " ")
		os.WriteFile(location, b, 0755)
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func main() {

	current_dir := path.Dir(os.Args[0])
	config_path := path.Join(current_dir, "config.json")
	l, e := os.Create(path.Join(current_dir, "restshut.log"))
	if e != nil {
		log.Fatal(e.Error())
	}
	log.SetOutput(l)
	defer l.Close()
	var err error
	c, err = ReadConfig(config_path)
	if err != nil {
		log.Fatalf("Could not read config, %s", err)
	}
	r := http.NewServeMux()
	r.HandleFunc("/", shutdownHandler)
	log.Printf("Listening on :%d", c.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r))
}
