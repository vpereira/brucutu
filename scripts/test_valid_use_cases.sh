#!/bin/bash
BRUCUTU=./build/brucutu

$BRUCUTU -u ssh://ssh -a 2222 -l root -p superpassword || exit -1
$BRUCUTU -u ssh://ssh -a 2222 -L samples/users.txt -P samples/passwd.txt || exit -1
$BRUCUTU -u pop3://email -l foo -p thepassword || exit -1
$BRUCUTU -u pop3://email -L samples/users.txt -P samples/passwd.txt || exit -1

exit 0
