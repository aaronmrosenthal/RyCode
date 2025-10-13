#!/bin/bash

# Test script to capture TUI screenshots for visual inspection

echo "Testing RyCode TUI cursor and layout..."

# Create screenshots directory
mkdir -p screenshots

# Test 1: Launch with splash, capture after splash
echo "Test 1: Splash screen and auto-session creation"
timeout 6 script -q screenshots/test1.txt bash -c "rycode" &
sleep 5
pkill -f "rycode"

# Test 2: Launch without splash
echo "Test 2: No splash, direct to chat"
timeout 3 script -q screenshots/test2.txt bash -c "rycode --no-splash" &
sleep 2
pkill -f "rycode"

# Test 3: Type some text and check cursor
echo "Test 3: Cursor positioning with text"
(
  sleep 0.5
  echo "hello world"
  sleep 1
) | timeout 3 script -q screenshots/test3.txt bash -c "rycode --no-splash"

echo "Screenshots saved to screenshots/"
echo "Review the .txt files to verify:"
echo "  - Cursor is solid (not blinking in static capture)"
echo "  - Cursor is positioned immediately after text"
echo "  - Input box is compact (minimal vertical space)"
echo "  - Splash lasts 4.5 seconds with proper fading"
