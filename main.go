package main

import (
	"fmt"
	"github.com/carmel/gooxml/document"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	logFile, err := os.Create("D:\\log.txt")
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

}

func main() {
	var dir string
	fmt.Println("请将所有Word文件放到同一文件夹中")
	fmt.Println("输入文件路径（只支持后缀名为.docx的Word文件): ")
	fmt.Scanln(&dir)

	start := time.Now()
	var wg sync.WaitGroup
	fileNames := make(chan string)
	//dir := "D:/data"
	go walkDir(dir, fileNames)

	for fileName := range fileNames {
		wg.Add(1)
		go func(file string) {
			countPages(file)
			wg.Done()
		}(fileName)
	}

	wg.Wait()
	timeUse := time.Since(start)
	fmt.Println("\n耗时: ", timeUse)
	abort := make(chan struct{})
	fmt.Println("按任意键关闭窗口")
	os.Stdin.Read(make([]byte, 1))
	abort <- struct{}{}
	select {
	case <-abort:
		os.Exit(1)
	}
}

func walkDir(dir string, fileNames chan<- string) {
	for _, entry := range dirents(dir) {
		fileNames <- path.Join(dir, entry.Name())
	}
	close(fileNames)
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

func countPages(file string) {
	doc, err := document.Open(file)
	if err != nil {
		log.Printf("error opening document: %s", err)
		output := `
【the page of this %s file is %d 】
-----------------------------------------------------------------------------------------
`
		fmt.Printf(output, file, 0)
		return
	}
	page := doc.AppProperties.X().Pages
	fmt.Printf("the page of this %s file is %d\n", file, *page)
}
