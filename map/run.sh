#!/bin/bash

CURR_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
docker run -it --net=host --expose 8001 --rm -v "$CURR_DIR":/usr/src/app -w /usr/src/app node:16 "$@"
