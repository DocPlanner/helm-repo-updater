package yq

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

const validKey = ".student-name"
const invalidKey = "student-name"

type Student struct {
	Name string `yaml:"student-name"`
	Age  int8   `yaml:"student-age"`
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

func writeSimpleYamlInFile() (*string, error) {
	yamlData, err := createSimpleYaml()

	if err != nil {
		return nil, fmt.Errorf("Error while Marshaling. %v", err)
	}

	tmpFile, err := createTempFile("yq-test")
	if err != nil {
		return nil, err
	}
	// Close the file
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	tmpFileName := tmpFile.Name()
	err = ioutil.WriteFile(tmpFileName, yamlData, 0644)
	if err != nil {
		return nil, err
	}
	return &tmpFileName, nil
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestQueryFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	expectedResult := "Sagar"
	result, err := QueryFile(validKey, *yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	if !(cmp.Equal(result, expectedResult)) {
		log.Fatalf("result and expectedResult for query %s in file %s are not equal", validKey, *yamlFile)
	}
}

func TestQueryFileInNotExistentFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()
	incorrectYamlFile := *yamlFile + "ts"

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	_, err = QueryFile(validKey, incorrectYamlFile)

	expectedError := fmt.Errorf(
		"open %s: no such file or directory",
		incorrectYamlFile,
	)
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestQueryFileIncorrectFile(t *testing.T) {
	tmpFile, err := createTempFile("yq-test")

	if err != nil {
		log.Fatal(err)
	}

	tmpFileName := tmpFile.Name()

	if _, err = tmpFile.WriteString("Hello World"); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFileName)

	_, err = QueryFile(validKey, tmpFileName)

	expectedError := fmt.Errorf(
		"returned non singular result for yq expression: '%s'",
		validKey,
	)
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestInplaceApplyInvalidKey(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	value := "Saga"
	err = InplaceApply(invalidKey, value, *yamlFile)

	expectedError := fmt.Errorf(
		"key %s doesn't start with '.'",
		invalidKey,
	)
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestInplaceApplyInexistentFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

	incorrectYamlFile := *yamlFile + "ts"

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	value := "Saga"
	err = InplaceApply(validKey, value, incorrectYamlFile)

	expectedError := fmt.Errorf(
		"stat %s: no such file or directory",
		incorrectYamlFile,
	)
	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}

func TestInplaceApply(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

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

	if !(cmp.Equal(value, *keyValueAfterInPlaceApply)) {
		log.Fatalf("value for key %s is not the expected after change value to %s in the file %s", key, value, *yamlFile)
	}
}

func TestReadKey(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	expectedKeyValue := "Sagar"

	keyValue, err := ReadKey(validKey, *yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	if !(cmp.Equal(*keyValue, expectedKeyValue)) {
		log.Fatalf("obtained value %s for key %s and expected value %s are not equal", *keyValue, validKey, expectedKeyValue)
	}
}

func TestReadKeyInvalidKey(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	_, err = ReadKey(invalidKey, *yamlFile)

	expectedError := fmt.Errorf(
		"key %s doesn't start with '.'",
		invalidKey,
	)

	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}

func TestReadKeyInexistentFile(t *testing.T) {
	yamlFile, err := writeSimpleYamlInFile()

	incorrectYamlFile := *yamlFile + "ts"

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(*yamlFile)

	_, err = ReadKey(validKey, incorrectYamlFile)

	expectedError := fmt.Errorf(
		"open %s: no such file or directory",
		incorrectYamlFile,
	)

	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}
