FROM node:20-alpine

# Java wird für Firestore & Auth Emulator benötigt
RUN apk add --no-cache openjdk21-jre bash curl

# Firebase CLI
RUN npm install -g firebase-tools@latest

WORKDIR /firebase
RUN firebase setup:emulators:firestore && \
    firebase setup:emulators:ui
COPY start.sh /start.sh

RUN chmod +x /start.sh
EXPOSE 4000 8080 9099
CMD ["/start.sh"]