/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"container/list"
	"fmt"
	"github.com/authzed/spicedb/pkg/schemadsl/dslshape"
	"github.com/authzed/spicedb/pkg/schemadsl/input"
	"github.com/authzed/spicedb/pkg/schemadsl/parser"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"regexp"
	"sort"
	"spicedb-dsl-validator/cmd/flags"
	"spicedb-dsl-validator/cmd/util"
	"strings"
)

type sNode struct {
	nodeType   dslshape.NodeType
	properties map[string]interface{}
	children   map[string]*list.List
}

func createAstNode(_ input.Source, kind dslshape.NodeType) parser.AstNode {
	return &sNode{
		nodeType:   kind,
		properties: make(map[string]interface{}),
		children:   make(map[string]*list.List),
	}
}

func (tn *sNode) GetType() dslshape.NodeType {
	return tn.nodeType
}

func (tn *sNode) Connect(predicate string, other parser.AstNode) {
	if tn.children[predicate] == nil {
		tn.children[predicate] = list.New()
	}

	tn.children[predicate].PushBack(other)
}

func (tn *sNode) MustDecorate(property string, value string) parser.AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

func (tn *sNode) MustDecorateWithInt(property string, value int) parser.AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

func getParseTree(currentNode *sNode, indentation int) string {
	parseTree := ""
	parseTree = parseTree + strings.Repeat(" ", indentation)
	parseTree = parseTree + fmt.Sprintf("%v", currentNode.nodeType)
	parseTree = parseTree + "\n"

	keys := make([]string, 0)

	for key := range currentNode.properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	//NodeTypeError //error-message
	for _, key := range keys {
		parseTree = parseTree + strings.Repeat(" ", indentation+2)
		parseTree = parseTree + fmt.Sprintf("%s = %v", key, currentNode.properties[key])
		parseTree = parseTree + "\n"
	}

	keys = make([]string, 0)

	for key := range currentNode.children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := currentNode.children[key]
		parseTree = parseTree + fmt.Sprintf("%s%v =>", strings.Repeat(" ", indentation+2), key)
		parseTree = parseTree + "\n"

		for e := value.Front(); e != nil; e = e.Next() {
			parseTree = parseTree + getParseTree(e.Value.(*sNode), indentation+4)
		}
	}

	return parseTree
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse spicedb schema",
	Long:  `./spicedb-dsl-validator parse --file-path tests/broken.zed`,
	Run:   parse,
}

func parse(cmd *cobra.Command, args []string) {
	filepath := flags.MustGetString("file-path", cmd.Flags())
	verbose := flags.MustGetBool("verbose", cmd.Flags())
	var schemacontent string
	err := util.ReadFileValueString(filepath, &schemacontent)
	if err != nil {
		glog.Error(err)
	}
	root := parser.Parse(createAstNode, input.Source(""), schemacontent)
	parseTree := getParseTree((root).(*sNode), 0)
	found := strings.TrimSpace(parseTree)
	if verbose {
		fmt.Println(found)
	}
	if strings.Contains(found, "error-message") {
		errorMessage, err := extractErrorMessage(found)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// Print the extracted error message
		fmt.Println("Extracted Error Message:", errorMessage)
	} else {
		fmt.Println("Parsed correctly")
	}

}

func extractErrorMessage(inputString string) (string, error) {
	// Define a regular expression to capture the error message
	re := regexp.MustCompile(`error-message\s*=\s*(.*)`)

	// Find the first match
	matches := re.FindStringSubmatch(inputString)

	// Check if there is a match
	if len(matches) < 2 {
		return "", fmt.Errorf("No error message found")
	}

	// Extract and trim the error message
	errorMessage := strings.TrimSpace(matches[1])

	return errorMessage, nil
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.PersistentFlags().String("file-path", "", "zed schema file")
	parseCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}