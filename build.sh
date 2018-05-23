#!/usr/bin/env bash

go build -o ccov
chmod +x ccov
mv -f ccov ~/bin/ccov