package main

import (
	"os"
	"fmt"
	"flag"
	"util"
	"path/filepath"
	"log"
	"os/exec"
	"io/ioutil"
)

var cmakeBuild string

var ProjectName string
var MainFile string

func init() {
	flag.StringVar(&ProjectName, "p", "-", "-p project_name")
	flag.StringVar(&MainFile, "f", "main", "-f main.c/c++")
}

func main() {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			util.PrintlnError(r)
		}
	}()
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(pwd)
	if ProjectName == "-" {
		_, ProjectName = filepath.Split(pwd)
	}

	fmt.Println(cmakeBuild)
	build(pwd)

}

func build(pwd string) {
	var cmakeClean = fmt.Sprintf("cd %s && cmake --build ./cmake-build-debug --target clean -- -j 4", pwd)
	cmakeBuild = "cmake --build ./cmake-build-debug --target all -- -j 8"
	cmakeBuild = fmt.Sprintf("cd %s && %s", pwd, cmakeBuild)

	//fmt.Println(cmakeClean + "\n" + cmakeBuild)

	shellRunExec(cmakeClean)
	shellRunExec(cmakeBuild)

	var llvmCov = fmt.Sprintf("cd %s && cd cmake-build-debug/CMakeFiles/%v.dir", pwd, ProjectName) + " && " +
		fmt.Sprintf("llvm-cov gcov -f -b %v.gcda", MainFile) + " && " +
		"lcov " +
		"--directory . " +
		"--base-directory . " +
		"--gcov-tool ~/bin/llvm-gcov.sh " +
		"--capture -o cov.info" + " && " +
		"genhtml cov.info -o outcov" + " && " +
		fmt.Sprintf("google-chrome "+
			"--app=file://%v/cmake-build-debug/CMakeFiles/"+
			"%v.dir/outcov/index.html", pwd, ProjectName)

	shellRunExec(llvmCov)

}

func shellRunExec(cmdArgs string) {
	var cmd = exec.Command("sh", "-c", cmdArgs)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err2 := cmd.StderrPipe()
	if err2 != nil {
		log.Fatal(err2)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	out, _ := ioutil.ReadAll(stdout)
	fmt.Printf("%s\n", out)

	errout, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", errout)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
