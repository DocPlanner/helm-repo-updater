package yq

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/pkg/errors"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

const outputFormat = "yaml"

var completedSuccessfully = false

// readFile takes a filepath and returns the byte value of the data within
func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	if _, err = file.Read(bytes); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return bytes, nil

}

func getYamlNode(rawData []byte) (yaml.Node, error) {
	var parsedData yaml.Node
	if err := yaml.Unmarshal(rawData, &parsedData); err != nil {
		return parsedData, fmt.Errorf("Error parsing yaml: %w", err)
	}
	return parsedData, nil
}

func disableYqlibLogging() {
	// this is the logger used by yq, set it to warning to hide trace and debug data
	logging.SetLevel(logging.WARNING, "")
}

// QueryFile get result of apply query to yaml file
func QueryFile(expression, filePath string) (interface{}, error) {
	disableYqlibLogging()
	b, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	node, err := getYamlNode(b)
	if err != nil {
		return nil, err
	}

	list, err := yqlib.NewAllAtOnceEvaluator().EvaluateNodes(expression, &node)
	if err != nil {
		return nil, err
	}
	nodes := []*yqlib.CandidateNode{}
	for el := list.Front(); el != nil; el = el.Next() {
		n := el.Value.(*yqlib.CandidateNode)
		nodes = append(nodes, n)
	}
	// should only match a single node
	if len(nodes) != 1 {
		return nil, errors.Errorf("returned non singular result for yq expression: '%s'", expression)
	}

	var result interface{}
	if err = nodes[0].Node.Decode(&result); err != nil {
		return nil, fmt.Errorf("Error decoding yaml.Node: %w", err)
	}
	return result, nil
}

// InplaceApply applies the yq expression to the given file
func InplaceApply(key, value string, targetFile string) error {
	if !strings.HasPrefix(key, ".") {
		return fmt.Errorf("key %s doesn't start with '.'", key)
	}

	expression := fmt.Sprintf("%s=\"%s\"", key, value)
	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(targetFile)
	out, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		return err
	}
	// need to indirectly call the function so  that completedSuccessfully is
	// passed when we finish execution as opposed to now
	defer func() { writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully) }()

	format, err := yqlib.OutputFormatFromString(outputFormat)
	if err != nil {
		return err
	}

	printerWriter := yqlib.NewSinglePrinterWriter(out)
	printer := yqlib.NewPrinter(printerWriter, format, true, false, 2, true)
	streamEvaluator := yqlib.NewStreamEvaluator()
	targetFiles := []string{
		targetFile,
	}
	err = streamEvaluator.EvaluateFiles(expression, targetFiles, printer, true)
	completedSuccessfully = err == nil

	return err
}

// ReadKey reads the value of the given key from the given file
func ReadKey(key string, targetFile string) string {
	if !strings.HasPrefix(key, ".") {
		return ""
	}
	query, err := QueryFile(key, targetFile)
	if err != nil {
		return ""
	}
	str := fmt.Sprintf("%v", query)
	return strings.TrimSpace(str)
}
