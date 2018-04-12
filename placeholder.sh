#!/bin/bash
set -e
LENGTH=$(head -c 4 | perl -ne 'print unpack("L", $_)')
declare -gr REQUEST=$(head -c $LENGTH)
VERSION=3000000

# Write a number as little-endian binary
function writelen()
{
    printf
}

# Require a command to be present, and quit if it's not
function require()
{
    if ! `command -v "$1" >/dev/null`; then
        OUTPUT="{\n  \"status\": \"error\"\n  "version": $VERSION,\n  \"message\": \"Required dependency '$1' is missing\"\n  \"code\": 1\n}"
        LANG=C LC_ALL=C LENGTH=${#OUTPUT}
        echo -n 
        exit 1
    fi
}

# Echo to stderr with failure status code
function fail()
{
    echo "$@" >&2
    return 1
}

# trim leading and trailing whitespace
function trim()
{
    echo "$1" | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*$//'
}

# Build an error response message
function error()
{
    if [ -n "$1" ]; then
        ERROR="$1"
    else
        ERROR=$(cat)
    fi
    [ -n "$2" ] && CODE="$2" || CODE=1
    OUTPUT="$(jq -n '.status = "error"' | jq --arg message "$ERROR" --arg code "$CODE" '.message = $message | .code = ($code|tonumber)')"
    perl -e "print pack('L', ${#OUTPUT})"
    echo -n "$OUTPUT"
}

# Wrap a shell command and bail with a valid response message on error
function wrap()
{
    . <({ STDERR=$({ STDOUT=$("$@"); EXIT=$?; } 2>&1; declare -p EXIT >&2; declare -p STDOUT >&2); declare -p STDERR; } 2>&1; )
    if [ $EXIT -gt 0 ] || [ -n "$STDERR" ]; then
        if [ -n "$STDERR" ]; then
            echo -n "$STDERR" | error
        elif [ -n "$STDOUT" ]; then
            echo -n "$STDOUT" | error
        else
            error "$@" $EXIT
        fi
    else
        OUTPUT="$(jq -n --arg version $VERSION --arg response "$STDOUT" '.status = "ok" | .version = $version | .response = ($response | if .[:1] == "{" then ( . | fromjson) else . end)')"
        perl -e "print pack('L', ${#OUTPUT})"
        echo -n "$OUTPUT"
    fi
    return $EXIT
}

# jq wrapper around the request
function rq()
{
    jq -r "$@" <<< "$REQUEST"
}

# Supply per-store configuration
function configure()
{
    set -e

    declare -A STORES
    . <(rq '.settings.stores | to_entries | map("STORES[\(.key|@sh)]=\(.value.path|@sh)")[]')

    for STORE in "${!STORES[@]}"; do
        STOREPATH=$(echo "${STORES["$STORE"]}" | sed 's/^~/$HOME/' | envsubst)
        [ -d "$STOREPATH" ] || fail "Store directory for '$STORE' does not exist: $STOREPATH"

        if [ -f "$STOREPATH/.browserpass.json" ]; then
            STORES["$STORE"]=$(cat "$STOREPATH/.browserpass.json")
        else
            STORES["$STORE"]=
        fi
    done

    [ -n "$PASSWORD_STORE_DIR" ] || PASSWORD_STORE_DIR="~/.password-store"
    OUTPUT="$(jq -n --arg defaultPath "$PASSWORD_STORE_DIR" '.defaultPath = $defaultPath')"

    for STORE in "${!STORES[@]}"; do
        OUTPUT=$(jq --arg store "$STORE" --arg settings "${STORES[$STORE]}" '.storeSettings[$store] = $settings' <<< "$OUTPUT")
    done

    echo "$OUTPUT"
}

# List all available logins by store
function list()
{
    set -e

    declare -A STORES
    . <(rq '.settings.stores | to_entries | map("STORES[\(.key|@sh)]=\(.value.path|@sh)")[]')
    for STORE in "${!STORES[@]}"; do
        STOREPATH=$(echo "${STORES["$STORE"]}" | sed 's/^~/$HOME/' | envsubst)
        [ -d "$STOREPATH" ] || fail "Store directory for '$STORE' does not exist: $STOREPATH"

        STORES[$STORE]="$(find -L "$STOREPATH" -type f -name '*.gpg' -printf '%P\n' | jq -R -s 'split("\n") | map(select(length > 0)) ')"
    done

    OUTPUT='{}'
    for STORE in "${!STORES[@]}"; do
        OUTPUT=$(jq --arg store "$STORE" --arg files "${STORES[$STORE]}" '.files[$store] = ($files|fromjson)' <<< "$OUTPUT")
    done

    echo "$OUTPUT"
}

# Fetch the specified fields from a single login
function fetch()
{
    set -e

    STORE="$(rq .store)"
    FILE="$(rq .file)"

    # sanity-check variables
    [ -n "$STORE" ] || fail ".store is not set"
    [ -n "$FILE" ] || fail ".file is not set"

    # get file path
    STOREPATH="$(rq --arg store "$STORE" '.settings.stores[$store].path//empty' | sed 's/^~/$HOME/' | envsubst)"
    [ -n "$STOREPATH" ] || fail "Store path is empty"
    [ -d "$STOREPATH" ] || fail "Store directory for '$STORE' does not exist: $STOREPATH"
    FILEPATH="$STOREPATH/$FILE"

    # get file contents
    [ -f "$FILEPATH" ] || fail "Requested file does not exist: $STORE:$FILE"
    DATA="$(gpg -q --decrypt "$FILEPATH")"
    
    # build output
    echo "$(jq -n --arg data "$DATA" '.data = $data')"
}

function run()
{
    case "$(rq .action)" in
        configure)  configure; return $?;;
        list)       list; return $?;;
        fetch)      fetch; return $?;;
        *)          echo "Unknown action: $(rq .action)" >&2; return 1;;
    esac
}

# Ensure dependencies are present
require cat
require envsubst
require find
require gpg
require grep
require head
require jq
require perl
require sed
require tac

# Run the client
wrap run
