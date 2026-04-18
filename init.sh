#!/bin/bash

log_ok() {
	printf "\e[32m[%s] %s\e[0m\n" "$(date)" "$1"
}

log_info() {
	printf "\e[33m[%s] %s\e[0m\n" "$(date)" "$1"
}

exit_with_error() {
	printf "\e[31m%s\e[0m\n" "$1"
	exit 1
}

{ 
	log_info "Regenerating database.."
	utils/create_db
} || {
	exit_with_error "failed to create db"; 
}

{ 
	log_info "Clearing messages.."
	utils/clear_messages 
} || {
	exit_with_error "Failed to clear messages"; 
}

{ 
	log_info "Getting dependencies.."
	go get 
} || { 
	exit_with_error "failed to get go dependencies"; 
}

{ 
	log_info "Building.."
	go build . 
} || { 
	exit_with_error "failed to build golang project"; 
}

log_ok "Done. to launch just make use of ./schat"
