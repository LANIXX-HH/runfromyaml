package functions

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	RestOut http.ResponseWriter
	ReqOut  *http.Request
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

func PrintSwitch(ctype color.Attribute, _level string, _output string, cstring ...interface{}) {
	switch _output {
	case "rest":
		PrintRest(ctype, _level, cstring...)
	case "file":
		PrintFile(_level, cstring...)
	case "stdout":
		PrintColor(ctype, _level, cstring...)
	}
}

func PrintFile(_level string, cstring ...interface{}) {
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

func PrintRestHeader() {
	log := logrus.New()
	//file
	log.Formatter = new(logrus.JSONFormatter)                      //default
	log.Formatter.(*logrus.JSONFormatter).PrettyPrint = false      // pretty print
	log.Formatter.(*logrus.JSONFormatter).DisableTimestamp = true  // remove timestamp from test output
	log.Formatter.(*logrus.JSONFormatter).DisableHTMLEscape = true // remove timestamp from test output

	log.Out = RestOut
	log.WithFields(logrus.Fields{
		"uri":        ReqOut.RequestURI,
		"method":     ReqOut.Method,
		"host":       ReqOut.Host,
		"remoteaddr": ReqOut.RemoteAddr,
		"header":     ReqOut.Header,
	}).Info("Header Information")
}

func PrintRest(ctype color.Attribute, _level string, cstring ...interface{}) {

	//fmt.Fprintln(RestOut, strings.Trim(fmt.Sprint(append(cstring, "")), "[]"))

	mystring := color.New(ctype)
	mystring.Fprintln(RestOut, cstring...)

	//fmt.Fprintln(RestOut, cstring...)
	//RestOut.Write([]byte("bla"))
	// log := logrus.New()
	// //file
	// log.Formatter = new(logrus.JSONFormatter)                      //default
	// log.Formatter.(*logrus.JSONFormatter).PrettyPrint = false      // pretty print
	// log.Formatter.(*logrus.JSONFormatter).DisableTimestamp = true  // remove timestamp from test output
	// log.Formatter.(*logrus.JSONFormatter).DisableHTMLEscape = true // remove timestamp from test output

	// // log.Formatter = new(logrus.TextFormatter)
	// // log.Formatter.(*logrus.TextFormatter).DisableColors = false
	// // log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true
	// // log.Formatter.(*logrus.TextFormatter).EnvironmentOverrideColors = true

	// // log.Formatter.(*logrus.TextFormatter).DisableLevelTruncation = false
	// // log.Formatter.(*logrus.TextFormatter).DisableQuote = true
	// // log.Formatter.(*logrus.TextFormatter).DisableSorting = true

	// // log.Formatter.(*logrus.TextFormatter).ForceColors = true
	// // log.Formatter.(*logrus.TextFormatter).ForceQuote = true

	// log.Out = RestOut
	// switch _level {
	// case "info":
	// 	log.Info(cstring...)
	// case "warn":
	// 	log.Warn(cstring...)
	// case "error":
	// 	log.Error(cstring...)
	// case "debug":
	// 	log.Debug(cstring...)
	// case "trace":
	// 	log.Trace(cstring...)
	// case "fatal":
	// 	log.Fatal(cstring...)
	// case "panic":
	// 	log.Panic(cstring...)
	// }
}

func PrintColor(ctype color.Attribute, _level string, cstring ...interface{}) {
	mystring := color.New(ctype)
	mystring.Println(cstring...)
}
