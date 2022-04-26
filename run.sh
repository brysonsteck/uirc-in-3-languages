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

echo "Compiling uirc in C..."
cd ./c
make -j$(nproc)
cd ..
echo

echo "Compiling uirc in Go..."
cd ./go
echo "go build ."
go build .
cd ..
echo

echo "Compiling uirc in Rust..."
cd ./rust
echo "cargo build"
cargo build
echo "cargo build --release"
cargo build --release
cd ..
echo 

echo "#####################"
echo "      Execution      "
echo "#####################"
echo

mkdir out

echo "Executing uirc in C..."
cd ./c
{ time ./uirc ../imgs/*.jpg > /dev/null 2>&1 ; } 2> ../out/c.txt
cd ..
echo

echo "Executing uirc in Go..."
cd ./go
{ time ./uirc ../imgs/*.jpg > /dev/null 2>&1 ; } 2> ../out/go.txt
cd ..
echo

echo "Executing uirc in Rust..."
cd ./rust
{ time ./target/release/uirc-rust ../imgs/*.jpg > /dev/null 2>&1 ; } 2> ../out/rust.txt
cd ..
echo

echo "#####################"
echo "       Results       "
echo "#####################"
echo

echo "Results for uirc in C:"
cat ./out/c.txt
echo

echo "Results for uirc in Go:"
cat ./out/go.txt
echo 

echo "Results for uirc in Rust:"
cat ./out/rust.txt
echo

rm -r ./out

