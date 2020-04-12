// libevent package
package main

import (
	"flag"
	"fmt"
	"github.com/richardpct/pkgsrc"
	"log"
	"os"
	"os/exec"
	"path"
)

var destdir = flag.String("destdir", "", "directory installation")
var pkg pkgsrc.Pkg

const (
	name     = "libevent"
	vers     = "2.1.11-stable"
	ext      = "tar.gz"
	url      = "https://github.com/libevent/libevent/releases/download/release-" + vers
	hashType = "sha256"
	hash     = "a65bac6202ea8c5609fd5c7e480e6d25de467ea1917c08290c521752f147283d"
)

func checkArgs() error {
	if *destdir == "" {
		return fmt.Errorf("Argument destdir is missing")
	}
	return nil
}

func configure() {
	fmt.Println("Waiting while launching autogen ...")
	path := os.Getenv("PATH")
	cmd1 := exec.Command("./autogen.sh")
	cmd1.Env = append(os.Environ(), "PATH="+path+":"+*destdir+"/bin")
	if out, err := cmd1.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}

	fmt.Println("Waiting while launching configure ...")
	cmd2 := exec.Command("./configure",
		"--prefix="+*destdir)
	cmd2.Env = append(os.Environ(), "CPPFLAGS="+"-I"+*destdir+"/include")
	if out, err := cmd2.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func build() {
	fmt.Println("Waiting while compiling ...")
	cmd := exec.Command("make", "-j"+pkgsrc.Ncpu)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func install() {
	fmt.Println("Waiting while installing ...")
	cmd := exec.Command("make", "install")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func main() {
	flag.Parse()
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	pkg.Init(name, vers, ext, url, hashType, hash)
	pkg.CleanWorkdir()
	if !pkg.CheckSum() {
		pkg.DownloadPkg()
	}
	if !pkg.CheckSum() {
		log.Fatal("Package is corrupted")
	}

	pkg.Unpack()
	wdPkgName := path.Join(pkgsrc.Workdir, pkg.PkgName)
	if err := os.Chdir(wdPkgName); err != nil {
		log.Fatal(err)
	}
	configure()
	build()
	install()
}
