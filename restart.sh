#!/bin/bash

cd $(dirname $0)
SELF_PATH=$(pwd)

DATE_STR=`date +"%Y%m%d"`
CORE_DUMP=1
TARGET_BIN_NAME="${SELF_PATH}/REST-GO-demo"
OPT=$@

#
function check_core_dump()
{
    if [ $CORE_DUMP -eq 1 ];
    then
        ulimit -c unlimited
        ulimit -s unlimited
    fi
}

function start()
{
    NUM=`ps -ef | grep -w "$TARGET_BIN_NAME" | grep -v grep | wc -l`
    if [ $NUM -ge 1 ]; then
        echo "Warnning: ${TARGET_BIN_NAME} is running!"
        return 1
    fi

    cd $SELF_PATH
    check_core_dump
    mkdir -p logs
    nohup $TARGET_BIN_NAME $OPT 2>&1 >> log/log.$DATE_STR  &
    echo "`date` ${TARGET_BIN_NAME} started"
}

run_kill()
{
    NUM=`ps -ef | grep -w "$TARGET_BIN_NAME" | grep -v grep | wc -l`
    if [ $NUM -eq 0 ]; then
        echo "${TARGET_BIN_NAME} is NOT running"
        return 0
    fi

    PID=`ps -ef | grep -w "$TARGET_BIN_NAME" | grep -v grep | gawk '{print $2}'`
    kill -9 $PID
    echo "$TARGET_BIN_NAME killed"
}

if [ "$1" = "start" ]; then
    start
elif [ "$1" = "stop" ]; then
    run_kill
elif [ "$1" = "kill" ]; then
    run_kill
elif [ "$1" = "kill9" ]; then
    run_kill
else
    start
fi