#!/bin/bash
BRUCUTU=./build/brucutu

echo "Test invalid URL"
$BRUCUTU -u :foo &&  exit 1

echo "invalid protocol"
$BRUCUTU -u xxx://localhost && exit 1

echo "Invalid set of arguments"
$BRUCUTU -u pop3://localhost -L samples/users.txt -l foo -p bar && exit 1

echo "False user and password for http basic auth"
$BRUCUTU -u http://http_basic_auth -l foo -p XXXX && exit 1

exit 0