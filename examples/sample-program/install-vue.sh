#! /bin/bash

echo "Checking for a previous Vue install..."
INSTALL_ONCE="src public package.json"

for FILE in $INSTALL_ONCE; do
    if [ -a $FILE ]; then
        # echo "$FILE was found"
        POSSIBLE_INSTALL=1
    fi
done

if [ -n "$POSSIBLE_INSTALL" ]; then
    echo "Looks like you already have a javascript project in this directory"
    exit
fi
echo "no install found...; installing sample files."

echo "INSTALLING VUE..."
npm create vite@latest vue-temp -- --template vue
cd vue-temp
cp -r src public package.json ..
if [ -a index.html ]; then
  cp index.html ..
fi 
cd ..
rm -r vue-temp
npm install

