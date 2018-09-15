#!/bin/bash

KB_ENV=${KB_ENV:-default}

ks delete ${KB_ENV} -c kubebench-quickstarter-volume
ks delete ${KB_ENV} -c kubebench-quickstarter-service
