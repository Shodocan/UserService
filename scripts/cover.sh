#!/bin/sh

go tool cover 2> /dev/null;
if [ $$ -eq 3 ]
then
		go get -u golang.org/x/tools/cmd/cover;
fi
go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html