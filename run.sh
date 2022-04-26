#!/bin/sh
#
# uirc in three different languages
# Bryson Steck (@brysonsteck)
#
# Runs all three programs and displays the output and time of each
#

type make gcc go cargo rustc
ERROR="$?"

if [ $ERROR -ne 0 ]; then 
  echo
  echo "run: one or more of the above programs are not installed. the programs in question have lines that end in \"not found\". please make sure that these programs are installed through your package manager"
  exit 1
fi

clear
echo "#####################"
echo "     Compilation     "
echo "#####################"
echo

echo "C"
echo "--------"
cd ./c
time make -j$(nproc)
cd ..
echo
ls -lh ./c | awk '$9 == "uirc" {print "filesize of C binary: " $5}'
wc ./c/uirc.c | awk '{print "line count of C file: " $1}'
echo
echo "Press [ENTER] to continue..."
read

echo "Go"
echo "--------"
cd ./go
echo "go build ."
time go build .
cd ..
echo
ls -lh ./go | awk '$9 == "uirc" {print "filesize of Go binary: " $5}'
wc ./go/uirc.go | awk '{print "line count of Go file: " $1}'
echo
echo "Press [ENTER] to continue..."
read

echo "Rust"
echo "--------"
cd ./rust
echo "cargo build --release"
time cargo build --release
cd ..
echo 
ls -lh ./rust/target/release/uirc-rust | awk '$9 == "uirc" {print "filesize of Rust binary: " $5}'
wc ./rust/src/main.rs | awk '{print "line count of Rust file: " $1}'
echo
echo "Press [ENTER] to continue..."
read

echo "#####################"
echo "      Execution      "
echo "#####################"
echo

echo "C"
echo "--------"
cd ./c
time ./uirc ../imgs/*.jpg
cd ..
echo
echo "Press [ENTER] to continue..."
read

echo "Go"
echo "--------"
cd ./go
time ./uirc ../imgs/*.jpg
cd ..
echo
echo "Press [ENTER] to continue..."
read

echo "Rust"
echo "--------"
cd ./rust
time ./target/release/uirc-rust ../imgs/*.jpg
cd ..
echo
