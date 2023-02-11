package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"text/template"

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

func GoTemplate(mymap map[string]string, mytemplate string) string {
	var writer bytes.Buffer
	t, err := template.New("todos").Parse(mytemplate)

	if err != nil {
		panic(err)
	}
	err = t.Execute(&writer, mymap)
	if err != nil {
		panic(err)
	}
	return writer.String()
}

func OpenAI(apiKey string, model string, prompt string, cmdtype string) map[string][]interface{} {
	// Erstellen Sie eine neue Anfrage an die OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/engines/text-davinci-003/completions", nil)
	if err != nil {
		fmt.Printf("Error creating API request: %s\n", err)
	}

	// Setzen Sie den API-Schlüssel und das Modell
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Stellen Sie die Anfrage an das Modell
	reqBody := map[string]interface{}{
		"prompt":      prompt + ". show a " + cmdtype + " example. Please do not write explanations. Please just a suggestion as" + cmdtype + " code.",
		"max_tokens":  100,
		"temperature": 0,
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error encoding request body: %s\n", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonReq))

	// Senden Sie die Anfrage an die API
	httpClient := &http.Client{}
	res, err := httpClient.Do(req.WithContext(context.Background()))
	if err != nil {
		fmt.Printf("Error sending API request: %s\n", err)
	}
	defer res.Body.Close()

	// Verarbeiten Sie die Antwort des Modells
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading API response: %s\n", err)
	}
	var response map[string][]interface{}
	json.Unmarshal(resBody, &response)
	// if err := json.Unmarshal(resBody, &response); err != nil {
	// 	fmt.Printf("Error decoding API response: %s\n", err)
	// }
	return response
}

// The function first checks if the key exists in the yblock map and if it does, it trims the value associated with the key and converts it to a string.
// If the expandenv key exists in the yblock map and is set to true, the function then expands any environment variables in the string.
// Finally, the function splits the string into an array of strings and returns it. If the key does not exist, then the function returns nil.

func ExtractAndExpand(yblock map[interface{}]interface{}, key string) []string {
	if reflect.ValueOf(yblock[key]).IsValid() {
		values := strings.Trim(fmt.Sprint(yblock[key]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
		}
		return strings.Fields(values)
	}
	return nil
}

func PrintShellCommandsAsYaml(commands []string) map[string]interface{} {

	mymap := map[string]interface{}{
		"logging": []map[string]interface{}{
			{
				"level": "info",
			},
			{
				"output": "stdout",
			},
		},
		"env": []map[string]interface{}{
			{
				"key":   "TEST",
				"value": "foo",
			},
			{
				"key":   "BLA",
				"value": "TEST",
			},
		},
		"cmd": []map[string]interface{}{
			// {
			// 	"type": "shell",
			// 	"values": []string{
			// 		"ls",
			// 	},
			// },
		},
	}

	for key, command := range commands {
		fmt.Println(key, command)
		mymap["cmd"] = append(mymap["cmd"].([]map[string]interface{}), map[string]interface{}{
			"type": "shell",
			"values": []string{
				command,
			},
		})
	}

	return mymap
}
