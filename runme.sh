#!/bin/bash

while go build main.go > /dev/null;
do sleep 1;
done;
./main