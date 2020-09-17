package main

import (
	"flag"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func listCmd(pkg string) []string {
	args := []string{"list",
		"-f", "{{if not .Standard}}{{.ImportPath}}{{end}}",
		"-deps",
		"-test",
		pkg}
	output, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		fmt.Printf("%s", output)
		panic(err)
	}
	ret := strings.Split(string(output), "\n")
	// for _, s := range ret {
	// 	if strings.Contains(s, "client-go") {
	// 		fmt.Printf("package %s imported %s\n", pkg, s)
	// 	}
	// }
	return ret
}

func main() {
	flag.Parse()
	visited := make(map[string]struct{})
	var newPkgs []string
	args := flag.Args()
	seed := "k8s.io/kubernetes/pkg/routes"
	if len(args) != 0 {
		seed = args[0]
	}
	newPkgs = listCmd(seed)
	var i int
	for {
		i++
		var remainingNewPkgs []string
		for _, pkg := range newPkgs {
			parts := strings.Split(pkg, " ")
			pkg = parts[0]
			pkg = strings.TrimSuffix(pkg, "_test")
			pkg = strings.TrimSuffix(pkg, ".test")
			if pkg == "" {
				continue
			}
			if _, ok := visited[pkg]; ok {
				continue
			}
			remainingNewPkgs = append(remainingNewPkgs, pkg)
			visited[pkg] = struct{}{}
		}
		// fmt.Printf("=========================%d========================\n", i)
		// for _, pkg := range remainingNewPkgs {
		// 	fmt.Println(pkg)
		// }
		if len(remainingNewPkgs) == 0 {
			break
		}
		newPkgs = []string{}
		for _, pkg := range remainingNewPkgs {
			newPkgs = append(newPkgs, listCmd(pkg)...)
		}
	}
	var allDeps []string
	for pkg, _ := range visited {
		allDeps = append(allDeps, pkg)
	}
	sort.Strings(allDeps)
	for _, pkg := range allDeps {
		fmt.Println(pkg)
	}
}
