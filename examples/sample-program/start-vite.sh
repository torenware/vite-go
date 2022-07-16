#! /usr/bin/env bash

JS_DIR=$1
SCRIPT=$0
VITE_PID=${VITE_PID:=/tmp/vite-script.pid}

usage() {
  echo "Usage: $SCRIPT JS_DIR" 2>/dev/null
}

if [ -z "$JS_DIR" ]; then
  usage
  exit 1
fi

if [ ! -d $JS_DIR ]; then
  echo "Script directory not found at $JS_DIR" 2>/dev/null
  exit 1
fi

# Make sure install has been run
if [ ! -d "$JS_DIR/node_modules" ]; then
  echo "New js dir; running npm install"
  BASE_DIR=$(pwd)
  cd $JS_DIR
  npm install
  cd $BASE_DIR
fi

cd $JS_DIR
node_modules/.bin/vite -l silent </dev/null &>/dev/null &
echo $! >$VITE_PID
disown
