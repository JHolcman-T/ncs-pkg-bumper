package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

type bumpType = string

const (
	major bumpType = "major"
	minor bumpType = "minor"
	patch bumpType = "patch"
)

func bumpVersion(source string, bump bumpType) string {
	var versionParts = strings.Split(source, ".")
	var partsCount = len(versionParts)
	if partsCount < 3 {
		for i := partsCount; i < 3; i++ {
			versionParts = append(versionParts, "0")
		}
	}
	if partsCount > 3 {
		versionParts = versionParts[:3]
	}

	var intVersionParts []int32
	for _, part := range versionParts {
		var value, err = strconv.ParseInt(part, 10, 32)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(errors.New("could not parse the version! Not all version parts are numeric!"))
			os.Exit(1)
		}
		intVersionParts = append(intVersionParts, int32(value))
	}

	switch bump {
	case major:
		return fmt.Sprintf("%d.0.0", intVersionParts[0]+1)
	case minor:
		return fmt.Sprintf("%d.%d.0", intVersionParts[0], intVersionParts[1]+1)
	case patch:
		return fmt.Sprintf("%d.%d.%d", intVersionParts[0], intVersionParts[1], intVersionParts[2]+1)
	}
	panic("Unknown bumpType")
}

var (
	filePath string
	bumpMode string
)

func init() {
	flag.StringVar(&filePath, "p", "", "Path to file. Defaults to './package-meta-data.xml'")
	flag.StringVar(&bumpMode, "m", "", "Mode of bump one of: ['major', 'minor', 'patch']")
	flag.Parse()
}

func main() {
	// validate args
	if filePath == "" {
		filePath = "./package-meta-data.xml"
	} else {
		pathInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if pathInfo.IsDir() {
			filePath = path.Join(filePath, "package-meta-data.xml")
		}
	}
	if bumpMode == "" {
		fmt.Println(errors.New("'-m' arg is missing"))
		os.Exit(1)
	}
	if bumpMode != major && bumpMode != minor && bumpMode != patch {
		fmt.Println(errors.New("'-m' flag has an unknown value"))
		os.Exit(1)
	}

	// Open our xmlFile
	xmlFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := io.ReadAll(xmlFile)

	// close file
	xmlFile.Close()

	// create document tree from file
	xmlDoc := etree.NewDocument()
	if err := xmlDoc.ReadFromBytes(byteValue); err != nil {
		panic(err)
	}

	// get package-name
	pkgName := xmlDoc.FindElement("/ncs-package/name")
	if pkgName == nil {
		fmt.Println(errors.New("could not grab 'name' tag from meta-data.xml"))
		os.Exit(1)
	}

	// get version tag
	pkgVersion := xmlDoc.FindElement("/ncs-package/package-version")
	if pkgVersion == nil {
		fmt.Println(errors.New("could not grab 'package-version' tag from meta-data.xml"))
		os.Exit(1)
	}

	// log package name
	fmt.Printf("--- Package: %s ---\n", pkgName.Text())

	// get version string
	currentVersion := pkgVersion.Text()

	// modify version
	newVersion := bumpVersion(currentVersion, bumpMode)

	// log versions
	fmt.Printf("Current version: %s\n", currentVersion)
	pkgVersion.SetText(newVersion)
	fmt.Printf("New version: %s\n", newVersion)

	// save to file
	xmlDoc.WriteToFile(filePath)
}
