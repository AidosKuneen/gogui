#!bin/sh
tsc test.ts
browserify test.js -o bundle.js