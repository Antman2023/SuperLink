package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	defer func() {
		fmt.Println("Press enter to exit...")
		fmt.Scanln()
	}()

	scriptDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Printf("Working directory: %s\n", scriptDir)

	if len(os.Args) < 2 {
		fmt.Println("Please drag & drop a directory to me.")
		return
	}

	targetStat, err := os.Stat(os.Args[1])

	if err != nil {
		fmt.Println("Check target error: ", err)
		return
	}

	if !targetStat.IsDir() {
		fmt.Println("Target is not a directory.")
		return
	}

	fmt.Println("About to process a directory: " + os.Args[1])

	var choice string
	fmt.Print("Continue running?(Y/N): ")
	fmt.Scanln(&choice)
	if strings.ToLower(choice) != "y" {
		return
	}

	modifiedParam := strings.Replace(os.Args[1], ":", "", -1)

	destinationPath := filepath.Join(scriptDir, modifiedParam)
	fmt.Printf(`Create dirctory: "%s"
	`, destinationPath)
	err = os.MkdirAll(destinationPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(`Moving directory: "%s" "%s"
	`, os.Args[1], destinationPath)
	err = copyDir(os.Args[1], destinationPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.RemoveAll(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Creating hard link...")
	cmd := exec.Command("cmd", "/c", "mklink", "/j", os.Args[1], destinationPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}
