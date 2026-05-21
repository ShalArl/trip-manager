#!/bin/sh

if [ -f "/firebase/data/firebase-export-metadata.json" ]; then
    echo "Importing from /firebase/data"
    exec firebase emulators:start \
        --project trip-manager-local \
        --only auth,firestore \
        --import=/firebase/data \
        --export-on-exit=/firebase/data
else
    echo "No data found, starting fresh"
    exec firebase emulators:start \
        --project trip-manager-local \
        --only auth,firestore \
        --export-on-exit=/firebase/data
fi