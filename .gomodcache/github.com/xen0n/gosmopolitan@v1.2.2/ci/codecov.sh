#!/bin/bash

set -e

curl -Os https://uploader.codecov.io/latest/linux/codecov
chmod +x codecov
./codecov
