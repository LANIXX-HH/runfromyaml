package functions

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
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

func PrintColor(ctype color.Attribute, cstring ...interface{}) {
	mystring := color.New(ctype)
	mystring.Println(cstring)
}
