package yq

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
	"gotest.tools/v3/assert"
)

const validKey = ".student-name"
const invalidKey = "student-name"
const contentSimpleFile = "Hello World"

type Student struct {
	Name string `yaml:"student-name"`
	Age  int8   `yaml:"student-age"`
}

func TestQueryFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	expectedResult := "Sagar"

	result, err := QueryFile(validKey, *yamlFile)

	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, result, expectedResult)
}

func TestQueryFileInNotExistentFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()
	incorrectYamlFile := *yamlFile + "ts"

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	_, err = QueryFile(validKey, incorrectYamlFile)

	expectedErrorMessage := fmt.Sprintf(
		"open %s: no such file or directory",
		incorrectYamlFile,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestQueryFileIncorrectFile(t *testing.T) {
	simpleTmpFile, err := writeSimpleTempFile()

	if err != nil {
		log.Fatal(err)
	}

	_, err = QueryFile(validKey, *simpleTmpFile)

	expectedErrorMessage := fmt.Sprintf(
		"returned non singular result for yq expression: '%s'",
		validKey,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestInplaceApplyInvalidKey(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	value := "Saga"
	err = InplaceApply(invalidKey, value, *yamlFile)

	expectedErrorMessage := fmt.Sprintf(
		`key %s doesn't start with '.'`,
		invalidKey,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestInplaceApplyInexistentFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	incorrectYamlFile := *yamlFile + "ts"

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	value := "Saga"
	err = InplaceApply(validKey, value, incorrectYamlFile)

	expectedErrorMessage := fmt.Sprintf(
		"stat %s: no such file or directory",
		incorrectYamlFile,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestInplaceApply(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	key := ".student-name"
	value := "Saga"

	err = InplaceApply(key, value, *yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	keyValueAfterInPlaceApply, err := ReadKey(key, *yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, value, *keyValueAfterInPlaceApply)
}

func TestReadKey(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	expectedKeyValue := "Sagar"

	keyValue, err := ReadKey(validKey, *yamlFile)

	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, *keyValue, expectedKeyValue)
}

func TestReadKeyInvalidKey(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	_, err = ReadKey(invalidKey, *yamlFile)

	expectedErrorMessage := fmt.Sprintf(
		"key %s doesn't start with '.'",
		invalidKey,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestReadKeyInexistentFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInTempFile()

	incorrectYamlFile := *yamlFile + "ts"

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	_, err = ReadKey(validKey, incorrectYamlFile)

	expectedErrorMessage := fmt.Sprintf(
		"open %s: no such file or directory",
		incorrectYamlFile,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestReadKeyIncorrectFile(t *testing.T) {
	simpleTmpFile, err := writeSimpleTempFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*simpleTmpFile)

	_, err = ReadKey(validKey, *simpleTmpFile)

	expectedErrorMessage := fmt.Sprintf(
		"returned non singular result for yq expression: '%s'",
		validKey,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func createTempFile(tempFilePrefix string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFilePrefix)
	if err != nil {
		return nil, fmt.Errorf("Cannot create temporary file. %v", err)
	}
	return tmpFile, nil
}

func createSimpleYaml() (out []byte, err error) {
	s1 := Student{
		Name: "Sagar",
		Age:  23,
	}

	yamlData, err := yaml.Marshal(&s1)
	if err != nil {
		return nil, err
	}
	return yamlData, nil
}

func writeSimpleYamlInTempFile() (*string, error) {
	yamlData, err := createSimpleYaml()

	if err != nil {
		return nil, fmt.Errorf("Error while Marshaling. %v", err)
	}

	tmpFile, err := createTempFile("yq-test")
	if err != nil {
		return nil, err
	}

	tmpFileName := tmpFile.Name()
	err = ioutil.WriteFile(tmpFileName, yamlData, 0644)
	if err != nil {
		return nil, err
	}
	// Close the file
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	return &tmpFileName, nil
}

func writeSimpleTempFile() (*string, error) {
	tmpFile, err := createTempFile("yq-test")

	if err != nil {
		return nil, err
	}

	tmpFileName := tmpFile.Name()

	if _, err = tmpFile.WriteString(contentSimpleFile); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	return &tmpFileName, nil
}
