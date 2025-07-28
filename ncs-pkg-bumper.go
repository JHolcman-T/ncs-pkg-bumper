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

type packageData struct {
	name        string
	version     string
	_versionTag *etree.Element
	path        string
	xmlDoc      *etree.Document
}

func strArray2IntArray(source []string) []int32 {
	// convert int strings to real ints
	var intVersionParts []int32
	for _, part := range source {
		var value, err = strconv.ParseInt(part, 10, 32)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(errors.New("could not parse the version! Not all version parts are numeric!"))
			os.Exit(1)
		}
		intVersionParts = append(intVersionParts, int32(value))
	}
	return intVersionParts
}

func parseVersion(source string) (out [3]int32) {
	// used for set-version input '-s' arg
	var versionParts = strings.Split(source, ".")
	if len(versionParts) != 3 {
		fmt.Println(errors.New("version must be a string of 3 numbers separated by dots"))
		os.Exit(1)
	}
	var intVersionParts = strArray2IntArray(versionParts)
	copy(out[:], intVersionParts)
	return out
}

func versionParts2String(source [3]int32) string {
	return fmt.Sprintf("%d.%d.%d", source[0], source[1], source[2])
}

func bumpVersion(source string, bump bumpType) string {
	// handle non standard input
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

	var intVersionParts = parseVersion(strings.Join(versionParts, "."))

	switch bump {
	case major:
		intVersionParts[0]++
	case minor:
		intVersionParts[1]++
	case patch:
		intVersionParts[2]++
	default:
		panic("Unknown bumpType")
	}

	return versionParts2String(intVersionParts)
}

var (
	filePath      string
	bumpMode      string
	setVersionArg string
)

func init() {
	flag.StringVar(&filePath, "p", "", "Path to file. Defaults to './package-meta-data.xml'")
	flag.StringVar(&bumpMode, "m", "", "Mode of bump one of: ['major', 'minor', 'patch']")
	flag.StringVar(&setVersionArg, "s", "", "Set version for the package. Syntax: <NUMBER>.<NUMBER>.<NUMBER> Example: 1.0.0")
	flag.Parse()
}

func readPackageData(path string) *packageData {
	// Open our xmlFile
	xmlFile, err := os.Open(path)
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

	return &packageData{
		name:        pkgName.Text(),
		version:     pkgVersion.Text(),
		_versionTag: pkgVersion,
		path:        path,
		xmlDoc:      xmlDoc,
	}
}

func setPackageVersion(version string, pkgData *packageData) {
	// log package name
	fmt.Printf("--- Package: %s ---\n", pkgData.name)

	// log versions
	fmt.Printf("Current version: %s\n", pkgData.version)
	pkgData._versionTag.SetText(version)
	fmt.Printf("New version: %s\n", version)

	// save to file
	pkgData.xmlDoc.WriteToFile(pkgData.path)
}

func opSetVersion() {
	parsedVersion := parseVersion(setVersionArg)
	stringifiedVersion := versionParts2String(parsedVersion)
	pkgData := readPackageData(filePath)
	setPackageVersion(stringifiedVersion, pkgData)
}

func opBumpVersion() {
	if bumpMode != major && bumpMode != minor && bumpMode != patch {
		fmt.Println(errors.New("'-m' flag has an unknown value"))
		os.Exit(1)
	}

	pkgData := readPackageData(filePath)
	bumpedVersion := bumpVersion(pkgData.version, bumpMode)
	setPackageVersion(bumpedVersion, pkgData)
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
	if bumpMode == "" && setVersionArg == "" {
		fmt.Println(errors.New("select mode of operation, set '-m' arg or use '-s' to set explicit version for the package"))
		os.Exit(1)
	}

	isSetVersionOp := setVersionArg != ""
	isBumpVersionOp := bumpMode != ""

	if isBumpVersionOp && isSetVersionOp {
		fmt.Println(errors.New("only one mode of operation is allowed"))
		os.Exit(1)
	}

	if isBumpVersionOp {
		opBumpVersion()
	}
	if isSetVersionOp {
		opSetVersion()
	}
}
