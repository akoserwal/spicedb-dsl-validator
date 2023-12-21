/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"container/list"
	"context"
	"fmt"
	"github.com/authzed/spicedb/pkg/schemadsl/compiler"
	"github.com/authzed/spicedb/pkg/schemadsl/dslshape"
	"github.com/authzed/spicedb/pkg/schemadsl/input"
	"github.com/authzed/spicedb/pkg/schemadsl/parser"
	"github.com/authzed/spicedb/pkg/schemautil"
	"github.com/spf13/cobra"
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
	parseonly := flags.MustGetBool("parse-only", cmd.Flags())

	var schemaContent string
	err := util.ReadFileValueString(filepath, &schemaContent)
	if err != nil {
		fmt.Println(err)
	}
	root := parser.Parse(createAstNode, input.Source(""), schemaContent)
	parseTree := getParseTree((root).(*sNode), 0)
	found := strings.TrimSpace(parseTree)
	if verbose {
		fmt.Println(found)
	}
	if parseonly {
		fmt.Println(found)
	} else {
		compiled, err := compiler.Compile(compiler.InputSchema{
			Source:       input.Source("schema"),
			SchemaString: schemaContent,
		}, new(string))
		if err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Complied error %s", err)
			//logger.Err(err).Msg("Complied error")
		}

		if compiled != nil {
			_, err := schemautil.ValidateSchemaChanges(context.Background(), compiled, false)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Schema validation error %s", err)
				//logger.Err(err).Msg("Schema validation error")
			}
		}
	}

}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.PersistentFlags().StringP("file-path", "f", "", "zed schema file path")
	parseCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	parseCmd.PersistentFlags().BoolP("parse-only", "p", false, "only parsed output")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
