#!/bin/bash
# LogAid Runner Script
# This script loads environment variables from .env file and runs LogAid

# Check if .env file exists
if [ -f ".env" ]; then
    # Export variables from .env file
    export $(grep -v '^#' .env | xargs)
    echo "Loaded environment variables from .env file"
else
    echo "Warning: .env file not found. Please create one based on .env.example"
    echo "You can also set environment variables manually:"
    echo "  export GEMINI_API_KEY=your_key_here"
    echo "  export OPENAI_API_KEY=your_key_here"
fi

# Run LogAid with all arguments passed through
./logaid "$@"
