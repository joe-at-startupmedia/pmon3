PROCESS_JSON="$1"
PROCESS_ID=$(echo "${PROCESS_JSON}" | jq '.id')
PROCESS_NAME=$(echo "${PROCESS_JSON}" | jq '.name')
echo "process event: ${PROCESS_ID} - ${PROCESS_NAME}"
