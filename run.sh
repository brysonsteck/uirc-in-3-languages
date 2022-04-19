#!/bin/sh
#
# uirc in three different languages
# Bryson Steck (@brysonsteck)
#
# Runs all three programs and displays the output and time of each
#

clear
echo "#####################"
echo "     Compilation     "
echo "#####################"
echo

echo "Compiling uirc in C..."
cd ./c
make
cd ..
echo

echo "Compiling uirc in Go..."
cd ./go
# compile in go
cd ..
echo

echo "Compiling uirc in Rust..."
cd ./rust
# compile in rust
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
# execute in go
cd ..
echo

echo "Executing uirc in Rust..."
cd ./rust
# execute in rust
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

