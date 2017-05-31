#!/bin/bash

ROOT=`pwd`

cd ${ROOT}
cd assets/scripts

ANGULAR=`find . -name 'angular.min.js'`
VENDOR=`find vendor -name '*.js' ! -name 'angular.min.js'`
ALL=`find . -name '*.js' ! -name 'app.js'  | grep -v "vendor"`
APP=`find . -name 'app.js'`
SCRIPTS=""

cd ${ROOT}

script_line () {
    FILE=`echo $1| sed -e "s/\.\///"`
    SCRIPT=`echo "<script src=\"/scoreboard/scripts/${FILE}\"></script>" | sed -e s"/\//∑/g"`
    SCRIPTS="${SCRIPTS}${SCRIPT}"
}


for LINE in ${ANGULAR}
do
    script_line $LINE
done

for LINE in ${VENDOR}
do
    script_line $LINE
done

for LINE in ${APP}
do
    script_line $LINE
done

for LINE in ${ALL}
do
    script_line $LINE
done

cd assets
sed -e "s/%script%/${SCRIPTS}/g" index.htmlp | sed -e "s/∑/\//g" >index.html