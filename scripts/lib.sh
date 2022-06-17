#!/bin/bash

function error {
  echo "ERROR : $@"
}

function info {
  echo "INFO  : $@"
}

function fatal {
  error $@
  exit 1
}

function check_installed {
  which $1 > /dev/null 2>&1 || fatal "$1 must be installed"
}

function req {
  [ -z $API_KEY ] && fatal "API_KEY is missing"
  [ -z $ENDPOINT ] && fatal "ENDPOINT is missing"

  METH=GET
  [ -z $2 ] || { METH=$1; shift; }
  URL=$1
  shift

  echo $@
  exit
  curl -s -X $METH -H "Authorization: basic $API_KEY" $@ $ENDPOINT/api/v1/$URL
}

function yn {
  echo "$@ (y/N)"
  read R
  [ "$R" == "y" ] && return 0 || return 1
}
