#!/bin/bash

echo "======================================="
echo "Go Todo API - Microservices Setup"
echo "======================================="
echo ""

# Check if .env file exists
if [ -f ".env" ]; then
    echo "✓ .env file already exists"
else
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "✓ .env file created"
    echo ""
    echo "⚠️  IMPORTANT: Please update the .env file with your Telegram credentials:"
    echo "   1. Open Telegram and search for @BotFather"
    echo "   2. Send /newbot command and follow instructions"
    echo "   3. Copy your bot token to TELEGRAM_BOT_TOKEN in .env"
    echo "   4. Search for @userinfobot on Telegram"
    echo "   5. Send any message and copy your ID to TELEGRAM_CHAT_ID in .env"
    echo ""
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi
echo "✓ Docker is running"

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install it and try again."
    exit 1
fi
echo "✓ docker-compose is available"

echo ""
echo "======================================="
echo "Setup completed!"
echo "======================================="
echo ""
echo "Next steps:"
echo "  1. Edit .env file with your Telegram credentials"
echo "  2. Run: make up"
echo "  3. Check logs: make logs"
echo "  4. Test API: curl http://localhost:8080/health"
echo ""
