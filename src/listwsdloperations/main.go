// lists operations within a WSDL
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const wsdldir = "../../wsdl"

// WSDL is an SOAP WSDL
type WSDL struct {
	Types     []interface{}   `xml:"schema"`
	Messages  []interface{}   `xml:"message"`
	PortTypes []OperationInfo `xml:"portType"`
}

// OperationInfo is a PortType
type OperationInfo struct {
	Name          string      `xml:"name,attr"`
	Documentation string      `xml:"documentation"`
	Operations    []Operation `xml:"operation"`
}

// Operation is an Operation under PortTypes
type Operation struct {
	Name          string `xml:"name,attr"`
	Documentation string `xml:"documentation`
}

func main() {

	// for all the wsdls in dir
	// read <wsdl:portType>
	// output <wsdl:documentation>
	// list all operations <wsdl:operation name>

	files, err := ioutil.ReadDir(wsdldir)
	if err != nil {
		log.Fatalln(err)
	}
	var wsdls int
	var ops int
	for _, f := range files {
		if strings.Contains(f.Name(), ".wsdl") {
			//fmt.Println(f.Name())
			desc, operations, err := listOperations(fmt.Sprintf("%s/%s", wsdldir, f.Name()))
			if err != nil {
				log.Printf("Can't read %s: %s\n", f.Name(), err)
			}
			wsdls++
			ops += len(operations)
			fmt.Printf("\"%s\",\"%s\",%v\n", f.Name(), desc, len(operations))
			
		}
	}
	fmt.Printf("%v wsdls, %v operations (%.2v)\n", wsdls, ops, (ops / wsdls))
}

func listOperations(xmlfile string) (string, []string, error) {
	var operations []string
	var desc string
	x, err := os.Open(xmlfile)
	if err != nil {
		return desc, operations, err
	}
	defer x.Close()

	b, _ := ioutil.ReadAll(x)

	var wsdl WSDL
	xml.Unmarshal(b, &wsdl)
	for _, p := range wsdl.PortTypes {
		desc = p.Documentation
		//fmt.Printf("\t%s\n", p.Documentation)
		for _, o := range p.Operations {
			//fmt.Printf("\t %s\n", o.Name)
			operations = append(operations,o.Name)
		}
	}
	return desc, operations, nil
}
