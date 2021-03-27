#!/bin/bash

go build main.go &
echo "build done" &
sudo ./main &