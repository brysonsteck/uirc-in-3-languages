/*

uirc-go: an unnecessary image ratio calculator, written in Go
Created by Bryson Steck (@brysonsteck on GitHub)

PROBLEM:
  * Currently can only read jpeg and png images.
The original C program (https://github.com/brysonsteck/uirc) can read all common image types. This program is a proof of concept.

Free and Open Source Software under the BSD 2-Clause License

BSD 2-Clause License

Copyright (c) 2022, Bryson Steck
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

package main

import (
  "os"
  "fmt"
  "net/http"
  "io"
  "image"
  _ "image/jpeg"
  _ "image/png"
  "strings"
)

const VERSION string = "0.1.0";
var rFlag bool = false;

func getBcf(width int, height int) (bcf int) {
  for i := 1; i <= width; i++ {
    for j := 1; j <= height; j++ {
      if width % i == 0 {
        if height % j == 0 && i == j {
          bcf = j;
        }
      }
    }
  }
  return bcf;
}

func readFile(file string, rFlag bool, req bool, url string) int {
  var width, height, factor int;
  var wuneven, huneven float32;
  imgFile, err := os.Open(file);
  if err != nil {
    fmt.Printf("uirc-go: %s: %s\n", file, err);
    os.Exit(6);
  }
  defer imgFile.Close();

  img, _, err := image.DecodeConfig(imgFile);
  if err != nil {
    fmt.Printf("uirc-go: %s: %s\n", file, err);
    os.Exit(3);
  }
  
  width = img.Width;
  height = img.Height;
  factor = getBcf(width, height);
  wuneven = float32(height) / float32(width);
  huneven = float32(width) / float32(height);

  if req {
    urlFile := strings.Split(url, "/")
    file = urlFile[len(urlFile) - 1]
  }

  if factor == 1 {
    if width < height {
      fmt.Printf("%s > 1:%.2f (uneven)", file, wuneven);
    } else {
      fmt.Printf("%s > %.2f:1 (uneven)", file, huneven);
    }
  } else {
    fmt.Printf("%s > %d:%d", file, width / factor, height / factor);
  }
  if rFlag {
    fmt.Printf(" [%dx%d]\n", width, height);
  } else {
    fmt.Println();
  }
  return 0;
}

func download(url string) int {
  response, err := http.Get(url)
	if err != nil {
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
    fmt.Printf("FAIL\nuirc-go: %s: request failed with code %d, trying local fs instead\n", url, response.StatusCode);
    return 4;
	}

	file, err := os.Create("/tmp/uirc.tmp")
	if err != nil {
    fmt.Printf("FAIL\nuirc-go: request complete, but cannot write file to /tmp for evaluation\n");
    os.Exit(9);
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
    fmt.Printf("uirc-go: error saving file\n");
    os.Exit(9);
	}
  fmt.Printf("ok\n");

  return 0;
}

func handleArg(arg string) {
  var complete int;
  var first, firstTwo, firstFour string;
  const help string = "USAGE: uirc [OPTIONS] IMAGE1 [IMAGE2] [...]\n\n" +

          "OPTIONS:\n" +
          "informational:\n" +
          "  -h, --help   \t: Display this message\n" +
          "  -l, --license\t: Display the license disclaimer for uirc-go (BSD 2-Clause)\n" +
          "  -v, --version\t: Display the version of uirc-go\n\n" +

          "functional:\n" +
          "  -r, --res    \t: Display the resolution of the image (in addition to the ratio)\n\n" +

          "HELP:\n" +
          "For more information on how to use uirc-go, open the man page uirc(1).\n";

  first = arg[0:1]; 

  if len(arg) >= 4 {
    firstFour = arg[0:4];
  }
  
  if len(arg) >= 2 {
    firstTwo = arg[0:2];
  }

  if "--" == firstTwo || "-" == first {
    if "--help" == arg || "-h" == arg {
      fmt.Printf("an unneccessary image ratio calculator, in Go! (uirc-go) v%s\n\n", VERSION);
      fmt.Printf("Copyright 2022 Bryson Steck\nFree and Open Source under the BSD 2-Clause License\n\n");
      fmt.Printf("%s\n", help);
      os.Exit(1);
    } else if "--license" == arg || "-l" == arg {
      fmt.Printf("uirc-go is Free and Open Source Software under the BSD 2-Clause License.\n");
      fmt.Printf("Please read the license regarding copying and distributing uirc.\n");
      fmt.Printf("https://github.com/brysonsteck/uirc/blob/master/LICENSE\n");
      os.Exit(1);
    } else if "--res" == arg || "-r" == arg {
      rFlag = true;
      return;
    } else if "--version" == arg || "-v" == arg {
      fmt.Printf("uirc-go v%s\n", VERSION);
      os.Exit(1);
    } else {
      fmt.Printf("uirc-go: invalid argument \"%s\"\nType \"uirc --help\" for help with arguments.\n", arg);
      os.Exit(5);
    }
  }

  if "http" == firstFour {
    fmt.Printf("downloading \"%s\"...", arg);
    download(arg);
    complete = readFile("/tmp/uirc.tmp", rFlag, true, arg);
    if complete != 0 {
      readFile(arg, rFlag, false, "");
    } else {
      os.Remove("/tmp/uirc.tmp");
    }
  } else {
    // if no more flags, run ratio calculations
    readFile(arg, rFlag, false, "");
  }
}

func main() {
  var runs int;

  if len(os.Args[1:]) == 0 {
    fmt.Println("uirc-go: at least one argument is required");
    os.Exit(1);
  }

  for _, a := range os.Args[1:] {
    handleArg(a);
    runs++;
  }

  if runs < 2 && rFlag == true {
    fmt.Println("uirc-go: at least one file/url is required");
    os.Exit(1);
  }
}
