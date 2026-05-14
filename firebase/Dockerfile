FROM node:20-alpine

# Java wird für Firestore & Auth Emulator benötigt
RUN apk add --no-cache openjdk21-jre bash curl

# Firebase CLI
RUN npm install -g firebase-tools@latest

WORKDIR /firebase

RUN firebase setup:emulators:firestore && \
    firebase setup:emulators:ui

EXPOSE 4000 8080 9099

CMD ["firebase", "emulators:start", \
"--project", "trip-manager-local", \
"--only", "auth,firestore", \
"--import=/firebase/data/latest", \
"--export-on-exit=/firebase/data/latest"]