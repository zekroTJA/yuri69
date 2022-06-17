#!/bin/bash

ENDPOINT=$1
API_KEY=$2

source $(dirname $(realpath $0))/lib.sh

[ -z $ENDPOINT ] || [ -z $API_KEY ] && \
  fatal "Usage: $(basename $0) [ENDPOINT] [API_KEY]"

check_installed curl
check_installed jq

SOUNDS=( $(req sounds | jq -r '.[] | .uid') )
N=${#SOUNDS[@]}

yn "Do you really want to DELETE ALL $N sounds on $ENDPOINT?" || exit

I=0

for SOUND in ${SOUNDS[@]}; do
  req DELETE sounds/$SOUND
  I=$(( I  + 1 ))
  info "Deleted Sound $I/$N"
done
