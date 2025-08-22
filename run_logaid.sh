#!/bin/bash
# LogAid Runner Script
# This script exports the necessary environment variables and runs LogAid

export GEMINI_API_KEY=AIzaSyAkA_IdOZZCiT9NQ85ccmzVKMU-jB2-QTc
export OPENAI_API_KEY=sk-proj-65wl5QEP4_j1JINr3RcFm7SJRKFbGcm33aOqKRGwVhjHV5A_9veBxd6y6cZYBkHN2DLAqJOx4T_BlbkFJ

# Run LogAid with all arguments passed through
./logaid "$@"
