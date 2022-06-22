#!/bin/sh
cd -- "$(dirname "$BASH_SOURCE")"
cd frontend && code . && cd ../server && code .
