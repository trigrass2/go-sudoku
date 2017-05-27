/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package sudokuparser

// FOO CPPFLAGS: -I/Users/jandersen/anaconda3/envs/opencvenv/include -I/Users/jandersen/anaconda3/envs/opencvenv/include/opencv2
// Foo CPPFLAGS: -I/usr/local/Cellar/opencv3/3.2.0/include -I/usr/local/Cellar/opencv3/3.2.0/include/opencv2

/*
#cgo CPPFLAGS: -I/usr/local/Cellar/opencv3/3.2.0/include -I/usr/local/Cellar/opencv3/3.2.0/include/opencv2
#cgo CXXFLAGS: --std=c++1z -stdlib=libc++
#cgo darwin LDFLAGS: -L/usr/local/Cellar/opencv3/3.2.0/lib -lopencv_core -lopencv_highgui -lopencv_imgcodecs -lopencv_imgproc -lopencv_ml -lopencv_objdetect -lopencv_photo
#include <stdlib.h>
#include "sudoku_parser.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var svmModelPath string

// Parse a Sudoku puzzle using a file path to a Sudoku image
func ParseSudokuFromFile(filename string) string {
	if !path.IsAbs(filename) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("Getwd failed: %s", err))
		}
		filename = path.Clean(path.Join(cwd, filename))
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return ParseSudokuFromByteArray(data)
}

// Parse a Sudoku puzzle from an image byte array
func ParseSudokuFromByteArray(data []byte) string {
	if svmModelPath == "" {
		svmModelPath = setupSVMModel()
		svmEnvVar := C.GoString(C.SVM_MODEL_VAR)
		err := os.Setenv(svmEnvVar, svmModelPath)
		if err != nil {
			panic(err)
		}

		fmt.Print(fmt.Sprintf("Set environment variable %s=%s\n", svmEnvVar, svmModelPath))
	}

	p := C.CBytes(data)
	defer C.free(p)

	parsed := C.ParseSudoku((*C.char)(p), C.int(len(data)), true)

	return C.GoString(parsed)
}

func setupSVMModel() string {
	tmpDir := path.Join(os.TempDir(), "sudokusolver")
	os.MkdirAll(tmpDir, os.ModePerm)

	tmpModelFile := path.Join(tmpDir, "sudokuSVMModel.yml")
	finfo, err := os.Stat(tmpModelFile)
	if err != nil {
		// create the file
		data, err := Asset("data/model4.yml")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(tmpModelFile, data, os.ModePerm)
		if err != nil {
			panic(err)
		}

		fmt.Println("SVM model file written to " + tmpModelFile)
	}

	if finfo.IsDir() {
		panic(tmpModelFile + " is a directory")
	}

	return tmpModelFile
}

// Parse a Sudoku puzzle from an image byte array
func TrainSudoku(trainConfigFile string) string {

	parsed := C.TrainSudoku(C.CString(trainConfigFile))

	return C.GoString(parsed)
}
