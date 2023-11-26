#!/bin/sh

# Environment variables USERNAME and PASSWORD must be set
htpasswd -cb /etc/nginx/.htpasswd foo bar
