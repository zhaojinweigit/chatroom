#!/bin/bash
gofmt -w .
go build -v server.go
go build -v cmdclient.go
