#!/bin/bash

go build main.go &
echo "build done" &
sleep 5 &
./main &