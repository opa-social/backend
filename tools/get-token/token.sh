#!/bin/bash
#
# Gets ID token for specific user from Firebase.
# Maintainer: Nick Pleatsikas <nick@pleatsikas.me>

getToken () {
    local token="$1"
    local username="$2"
    local password="$3"

    curl "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$token" \
        --silent \
        -H 'Content-Type: application/json' \
        --data-binary @<(cat <<EOF
{
    "email": "$username",
    "password": "$password",
    "returnSecureToken": true
}
EOF
) | jq .idToken
        
}

main () {
    # Storage variables for flags.
    local token
    local username
    local password

    while [[ ! $# -eq 0 ]]; do
        case "$1" in
            --username=*)
                IFS="=" read -ra USERNAME_FLAG <<< "$1"
                username="${USERNAME_FLAG[1]}"

                ;;
            --password=*)
                IFS="=" read -ra PASSWORD_FLAG <<< "$1"
                password="${PASSWORD_FLAG[1]}"

                ;;
            --token=*)
                IFS="=" read -ra TOKEN_FLAG <<< "$1"
                token="${TOKEN_FLAG[1]}"

                ;;
        esac
        shift
    done

    if [[ -z "$token" || -z "$username" || -z "$password" ]]; then
        echo -e "Token, username, or password is missing!"
        exit 1
    fi

    getToken "$token" "$username" "$password"

    exit 0
}

main "$@"