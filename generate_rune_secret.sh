#!/bin/bash

base64_data=$(openssl rand -base64 32)
echo -n "$base64_data" | base64 -d | xxd -p

