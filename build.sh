#!/usr/bin/env bash

go build -ldflags="-H windowsgui -linkmode internal -s -w"
