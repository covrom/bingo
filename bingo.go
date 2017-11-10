package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	var err error
	var data []byte
	var w *os.File
	var ishtml bool

	if len(os.Args) != 4 {
		log.Fatal("Syntax: filename_or_web_link package_name var_name")
	}
	fin := os.Args[1]
	pkgname := os.Args[2]
	varname := os.Args[3]

	if strings.HasPrefix(fin, "http://") || strings.HasPrefix(fin, "https://") {
		rsp, err := http.Get(fin)
		if err != nil {
			log.Fatal(err)
		}
		defer rsp.Body.Close()
		data, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal(err)
		}

		u, err := url.Parse(fin)
		if err != nil {
			log.Fatal(err)
		}

		pthb := path.Base(u.Path)

		ishtml = strings.HasSuffix(pthb, ".html") ||
			strings.HasSuffix(pthb, ".htm") ||
			strings.HasSuffix(pthb, ".css") ||
			strings.HasSuffix(pthb, ".js")

		fin = "./" + pthb + ".go"

	} else {
		data, err = ioutil.ReadFile(fin)
		if err != nil {
			log.Fatal(err)
		}

		ishtml = strings.HasSuffix(fin, ".html") ||
			strings.HasSuffix(fin, ".htm") ||
			strings.HasSuffix(fin, ".css") ||
			strings.HasSuffix(fin, ".js")

		fin = fin + ".go"
	}

	w, err = os.Create(fin)

	if err != nil {
		log.Fatal(err)
	}

	defer w.Close()

	fmt.Fprintln(w, "package "+pkgname)
	fmt.Fprintln(w, "// "+os.Args[1])

	if ishtml {

		fmt.Fprint(w, "var "+varname+" = []byte(`")
		sdata := string(data)
		for _, r := range sdata {
			if r != '`' {
				fmt.Fprint(w, string(r))
			} else {
				fmt.Fprint(w, "`+\"`\"+`")
			}
		}
		fmt.Fprintln(w, "`)")

	} else {

		fmt.Fprintln(w, "var "+varname+" = []byte{")
		for i, v := range data {
			if (v != 0x5c) && (v != 0x27) && ((v >= 0x20 && v <= 0x5f) || (v >= 0x61 && v <= 0x7e)) {
				fmt.Fprintf(w, "'%s',", string(v))
			} else {
				fmt.Fprintf(w, "%#x,", v)
			}
			if i%16 == 15 {
				fmt.Fprintln(w)
			}
		}
		if len(data)%16 != 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, "}")

	}
	w.Sync()
	log.Println("Done "+os.Args[1])
}
