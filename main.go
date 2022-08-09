package main

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("USAGE: translate_check <original_snippet_file.json> <translated_snippet_file.json>\n")
		fmt.Printf("\n")
		fmt.Printf("This tool will compare the two files and check if any keys are missing from the translated file\n")
		fmt.Printf("and if there are any keys in the translated file that is not in the original file\n")
		fmt.Printf("\n")
		os.Exit(1)
	}

	originalTranslateData, err := readTranslateJson(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading original file: %s\n", err)
		os.Exit(1)
	}

	translatedData, err := readTranslateJson(os.Args[2])
	if err != nil {
		fmt.Printf("Error reading translated file: %s\n", err)
		os.Exit(1)
	}

	allOrigKeys := extractAllKeys(originalTranslateData)
	allTranslatedKeys := extractAllKeys(translatedData)

	missingKeys := make([]string, 0)
	canBeDeleted := make([]string, 0)

	for _, key := range allOrigKeys {
		if !slices.Contains(allTranslatedKeys, key) {
			missingKeys = append(missingKeys, key)
		}
	}

	for _, key := range allTranslatedKeys {
		if !slices.Contains(allOrigKeys, key) {
			canBeDeleted = append(canBeDeleted, key)
		}
	}

	if len(missingKeys) > 0 {
		fmt.Println("These keys are missing:")
		for _, key := range missingKeys {
			fmt.Println(key)
		}
		fmt.Println("")
	} else {
		fmt.Println("No keys are missing")
	}

	if len(canBeDeleted) > 0 {
		fmt.Println("These keys can be deleted:")
		for _, key := range canBeDeleted {
			fmt.Println(key)
		}
		fmt.Println("")
	} else {
		fmt.Println("No keys needs to be deleted")
	}

}

func readTranslateJson(filename string) (map[string]interface{}, error) {
	jsonRawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var JsonParsedData map[string]interface{}

	err = json.Unmarshal(jsonRawData, &JsonParsedData)
	if err != nil {
		return nil, err
	}

	return JsonParsedData, nil
}

func extractAllKeys(data map[string]interface{}) []string {

	keys := make([]string, 0)

	for key, value := range data {

		nestedKeys := make([]string, 0)
		if valueA, ok := value.(map[string]interface{}); ok {
			nestedKeys = extractAllKeys(valueA)
		}

		if len(nestedKeys) > 0 {
			for _, nk := range nestedKeys {
				keys = append(keys, fmt.Sprintf("%s.%s", key, nk))
			}
		} else {
			keys = append(keys, key)
		}

	}
	return keys
}
