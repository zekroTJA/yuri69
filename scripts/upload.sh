#!/bin/bash

LOCATION=$1
ENDPOINT=$2
API_KEY=$3

[ -z $LOCATION ] || [ -z $ENDPOINT ] || [ -z $API_KEY ] && {
  echo "Usage: $(basename $0) [FILE_LOCATION] [ENDPOINT] [API_KEY]"
}

function check_installed {
  which $1 > /dev/null 2>&1 || {
    echo "ERROR : $1 must be installed"
    exit 1
  }
}

check_installed curl
check_installed jq

[ -d $LOCATION ] || {
  echo "ERROR : $LOCATION does not exist or is not a directory"
}

for FILE in $LOCATION/*; do
  NAME=${FILE%.*}
  BASENAME=$(basename $NAME)

  MIME=$(file $FILE -i | awk -F: '{print $2}')
  MIME=${MIME%;*}

  echo "INFO : Uploading file $FILE ..."
  RES=$(curl -Ls \
    -X PUT \
    -F "file=@${FILE};type=${MIME}" \
    -H "Authorization: basic ${API_KEY}" \
      ${ENDPOINT}/api/v1/sounds/upload)

  UPLOAD_ID=$(echo $RES | jq -r .upload_id)

  curl -Ls \
    -X POST \
    --data "{ \"upload_id\": \"${UPLOAD_ID}\", \"uid\": \"${BASENAME}\", \"normalize\": true }" \
    -H "Content-Type: application/json" \
    -H "Authorization: basic ${API_KEY}" \
      ${ENDPOINT}/api/v1/sounds/create

  echo "{ \"upload_id\": \"${UPLOAD_ID}\", \"uid\": \"${BASENAME}\", \"normalize\": true }"

  echo "INFO : Sound $BASENAME successfully uploaded"
done
