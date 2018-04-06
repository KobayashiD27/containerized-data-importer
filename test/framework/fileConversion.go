package framework

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kubevirt/containerized-data-importer/pkg/image"
)

var formatTable = map[string]func(string) (string, error){
	image.ExtGz:    transformGz,
	image.ExtXz:    transformXz,
	image.ExtTar:   transformTar,
	image.ExtQcow2: transformQcow2,
	"":             transformNoop,
}

func FormatTestData(srcFile string, targetFormats ...string) (string, error) {

	var (
		err     error
		outFile = srcFile
	)

	for _, tf := range targetFormats {
		if err != nil {
			break
		}
		i := 0
		for ftkey, ffunc := range formatTable {
			if tf == ftkey {
				outFile, err = ffunc(outFile)
				break
			} else if i == len(formatTable)-1 {
				err = fmt.Errorf("format extension %q not recognized\n", tf)
			}
			i ++
		}
	}

	if err != nil {
		return "", fmt.Errorf("FormatTestData: %v", err)
	}
	return outFile, nil
}

func transformFile(srcFile, ext, osCmd string, osArgs ...string) (string, error) {
	//osArgs = append(osArgs, srcFile)
	cmd := exec.Command(osCmd, osArgs...)
	cout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("transformFile: command erred: %v\n"+
			"transformFile: command output: %v", err, string(cout))
	}
	f := srcFile + ext
	_, err = os.Stat(f)
	if err != nil {
		return "", fmt.Errorf("transformFile: error stat-ing file: %v\n", err)
	}
	return f, nil
}

func transformTar(srcFile string) (string, error) {
	args := []string{"-cf", srcFile + image.ExtTar, srcFile}
	return transformFile(srcFile, image.ExtTar, "tar", args...)
}

func transformGz(srcFile string) (string, error) {
	return transformFile(srcFile, image.ExtGz, "gzip", srcFile)
}

func transformXz(srcFile string) (string, error) {
	return transformFile(srcFile, image.ExtXz, "xz", srcFile)
}

func transformQcow2(srcfile string) (string, error) {
	args := []string{"convert", "-f", "raw", "-O", "qcow2", srcfile, srcfile + image.ExtQcow2}
	return transformFile(srcfile, image.ExtQcow2, "qemu-img", args...)
}

func transformNoop(srcFile string) (string, error) {
	return srcFile, nil
}
