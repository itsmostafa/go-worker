.PHONY: run
include .env
export

# target: run - Run dev scheduler.
run:
	go run .
