package functions

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

//Check error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// WriteFile write a file
func WriteFile(file string, path string, perm os.FileMode) {
	bytefile := []byte(file)
	err := ioutil.WriteFile(os.ExpandEnv(path), bytefile, perm)
	Check(err)
}

//ReadFile read file
func ReadFile(file string) {
	content, err := ioutil.ReadFile(os.ExpandEnv(file))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)
}

//Remove element from slice
func Remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func Config() {

}

func PrintColor(ctype color.Attribute, _level string, _output string, cstring ...interface{}) {

	log := logrus.New()

	switch _output {
	case "stdout":
		// stdout
		log.Formatter = new(logrus.TextFormatter)
		log.Formatter.(*logrus.TextFormatter).DisableLevelTruncation = true
		log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
		log.Formatter.(*logrus.TextFormatter).ForceColors = true
		log.Formatter.(*logrus.TextFormatter).EnvironmentOverrideColors = true

		log.Out = os.Stdout
		mystring := color.New(ctype)
		mystring.Println(cstring...)
	case "file":
		// file
		log.Formatter = new(logrus.JSONFormatter)                      //default
		log.Formatter.(*logrus.JSONFormatter).PrettyPrint = true       // pretty print
		log.Formatter.(*logrus.JSONFormatter).DisableTimestamp = false // remove timestamp from test output
		dir := os.TempDir()
		file, err := os.OpenFile(dir+"logrus-"+time.Now().Format("20060102")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}

		log.Info(cstring...)
		switch _level {
		case "info":
			log.Info(cstring...)
		case "warn":
			log.Warn(cstring...)
		case "error":
			log.Error(cstring...)
		case "debug":
			log.Debug(cstring...)
		case "trace":
			log.Trace(cstring...)
		case "fatal":
			log.Fatal(cstring...)
		case "panic":
			log.Panic(cstring...)
		}
	}
}
