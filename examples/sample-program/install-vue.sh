#! /bin/bash

echo "INSTALLING VUE..."
npm create vite@latest vue-temp -- --template vue
cd vue-temp 
cp -r src public package.json ..
cd ..
rm -r vue-temp
npm install

