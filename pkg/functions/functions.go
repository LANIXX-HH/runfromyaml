package functions

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// WriteFile write a file
func WriteFile(file string, path string, perm os.FileMode) {
	bytefile := []byte(file)
	err := os.WriteFile(os.ExpandEnv(path), bytefile, perm)
	if err != nil {
		panic(err)
	}
}

//ReadFile read file
func ReadFile(file string) {
	content, err := os.ReadFile(os.ExpandEnv(file))
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

	switch _output {
	case "stdout":
		mystring := color.New(ctype)
		mystring.Println(cstring...)
	case "file":
		log := logrus.New()
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
