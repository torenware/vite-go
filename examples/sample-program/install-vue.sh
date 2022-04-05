#! /bin/bash

d_flag='' # dev
config_file='' # config

verbose='false'

print_usage() {
    printf "Usage: $0 [-d] [-c config_file] [install_dir]"
}

while getopts ':cd' flag; do
    case "${flag}" in
        d) d_flag='true' ;;
        c) config_file="${OPTARG}" ;;
        v) verbose='true' ;;
        *) print_usage
        exit 1 ;;
    esac
done

shift $(( OPTIND - 1 ))
inst_dir=${1-frontend}

if [ -n "$d_flag" ]; then
    echo dev mode
else
    echo prod mode
fi

if [ -z "$config_file"]; then
    config_file=vite.config.js
fi
echo "use config $config_file"
echo install_dir: $inst_dir

echo "Checking for a previous Vue install..."
INSTALL_ONCE=$inst_dir

for FILE in $INSTALL_ONCE; do
    if [ -a $FILE ]; then
        # echo "$FILE was found"
        POSSIBLE_INSTALL=1
    fi
done

if [ -n "$POSSIBLE_INSTALL" ]; then
    echo "Looks like you already have a javascript project in $inst_dir"
    exit
fi
echo "no install found...; installing sample files."

echo "INSTALLING VUE..."
npm create vite@latest -y $inst_dir -- --template vue
if [ -a $config_file ]; then
    cp $config_file $inst_dir
fi
cd $inst_dir
npm install
cd ..

