package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// Check if the filepath is valid and exists on the disk
func CheckFileExists(filePath string) bool {
	//Get the current directory
	//pwd, _ := os.Getwd()
	//Check if the file exists
	_, err := os.Stat(filePath)
	//Return the result
	return !os.IsNotExist(err)
}

// Read all lines in the file
func ReadLinesFromFile(filePath string) ([]string, error) {
	//Get the current directory
	//pwd, _ := os.Getwd()
	//Check if the file exists
	exists := CheckFileExists(filePath)
	if !exists {
		return nil, errors.New("file does not exist")
	}
	//Open the file
	file, err := os.Open(filePath)
	//Check if an error occured when opening the file
	if err != nil {
		return nil, err
	}
	//Close the file at the end of the function
	defer file.Close()

	//Read lines from the file and append it to returning splice
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	lines := []string{}
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	//Return the lines
	return lines, nil
}

type ExploitRequest struct {
	URL        string
	Method     string
	Headers    map[string][]string
	Cookies    map[string]string
	GetParams  map[string][]string
	PostParams map[string][]string
	Body       string
}

func (er *ExploitRequest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(er)
}

// Create python exploit code from request
func CreatePythonExploitCode(rawRequest string, exploitTemplatePath string) (string, error) {
	ioReader := bytes.NewReader([]byte(rawRequest))
	reader := bufio.NewReader(ioReader)
	req, err := http.ReadRequest(reader)
	//Check if an error occured when reading the raw request
	if err != nil {
		return "", errors.New("could not read the raw request into struct, " + err.Error())
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	exploitReq := ExploitRequest{}

	for _, cookie := range req.Cookies() {
		exploitReq.Cookies[cookie.Name] = cookie.Value
	}
	exploitReq.Method = req.Method
	exploitReq.Headers = req.Header
	exploitReq.Headers["Host"] = make([]string, 0)
	exploitReq.Headers["Host"] = append(exploitReq.Headers["Host"], req.Host)

	var protocol string = "http"
	if req.TLS != nil {
		protocol = "https"
	}
	exploitReq.URL = protocol + "://" + req.Host + req.URL.RawPath
	exploitReq.GetParams = req.URL.Query()
	err = req.ParseForm()
	if err != nil {
		return "", err
	}
	exploitReq.PostParams = req.PostForm

	//Read the data from the body
	bodyData, err := io.ReadAll(req.Body)
	//Check if the body could have been read
	if err != nil {
		return rawRequest, errors.New("could not read the request body")
	}

	//Reassign the body so other function can read the data
	req.Body = io.NopCloser(bytes.NewReader(bodyData))
	exploitReq.Body = string(bodyData)

	tmpl, err := template.New("exploit.tmpl").Funcs(funcMap).ParseFiles(exploitTemplatePath)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	var exploitCode bytes.Buffer
	err = tmpl.Execute(&exploitCode, &exploitReq)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return exploitCode.String(), nil
}
