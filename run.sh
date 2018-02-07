#!/bin/bash

set -e

echo "Building app ..."
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH=$DIR

cd src
glide install
cd ..

echo "glide install processed."

go install gomypartition

echo "App installed."

#exec ./bin/gomypartition wait

echo "Press [CTRL+C] to stop.."
while true
do
	sleep 1
done
