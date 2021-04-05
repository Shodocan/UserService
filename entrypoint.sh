#!/bin/ash
set -e

token=$(wget --header "Metadata-Flavor: Google" \
    "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token" \
    -O - | grep 'access_token' | sed -r 's/^[^:]*:(.*)$/\1/')

wget --header "authorization: Bearer $token" "https://secretmanager.googleapis.com/v1/projects/$PROJECT_ID/secrets/$SECRET/versions/$SECRET_VER:access" -O - \
    | grep 'data' | sed -r 's/^[^:]*:(.*)$/\1/' | sed -e 's/^ //' -e 's/^"//' -e 's/"$//' | base64 -d \
    > .env

./$1