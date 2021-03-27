#!/bin/bash

go build main.go &
echo "build done" &
lxterminal -e "./main"