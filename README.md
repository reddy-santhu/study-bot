# Personalized Study Assistant Bot

This Discord bot is designed to help students and learners organize their study sessions, track progress, and stay productive. It incorporates features like study goal management, Pomodoro timers, smart reminders, AI-powered study resources (doubt solving, summaries, flashcard generation), and the ability to upload and interact with PDF documents.

## Table of Contents

-   [Overview](#overview)
-   [Features](#features)
-   [Prerequisites](#prerequisites)
-   [Setup](#setup)
-   [Configuration](#configuration)
-   [Running the Bot](#running-the-bot)
-   [Command Reference](#command-reference)
-   [Contributing](#contributing)
-   [License](#license)

## Overview

The Personalized Study Assistant Bot is built using Golang and leverages the Discord API to provide a seamless study companion within your Discord server.  It enhances productivity and learning through various interactive features, including AI-powered assistance via the Gemini API, study session management, and persistent data storage using MongoDB.

## Features

-   **Study Goal Management:**
    -   `/study set <goal>`: Sets a new study goal.
    -   `/study list`: Lists your current study goals.
    -   `/study remove <goal number>`: Removes a specific study goal (by its number in the list).

-   **Pomodoro Timer:**
    -   `/pomo start <work time> <break time> <task>`: Starts a Pomodoro timer (work/break times in minutes).
    -   `/pomo stop`: Stops the running Pomodoro timer.
    -   `/pomo status`: Shows remaining time in the current Pomodoro session.

-   **Smart Study Reminders:**
    -   `/remind <time> <task>`: Sets a reminder for a task at a given time (HH:MM or YYYY-MM-DD HH:MM).

-   **AI-Powered Features (Gemini API):**
    -   `/ask question: <question>`: Asks Gemini a study-related question.
    -   `/summarize text: <long text>`: Gets a concise summary of a long text passage.

-   **PDF Interaction:**
    -   `/pdf`: Uploads a PDF file for later interaction (content stored in MongoDB).
    -   `/viewpdf`: Shows a numbered list of uploaded PDFs (for the user).
    -   `/askAi <pdf_number> <question>`: Asks Gemini a question related to the PDF's content, maintaining conversational context.

## Prerequisites

-   **Go:** Go programming language (version 1.18 or later recommended).  Install from [go.dev](https://go.dev/dl/).
-   **Discord Bot Token:** Obtain a bot token from the [Discord Developer Portal](https://discord.com/developers/applications).  **Enable the "Message Content Intent"**.
-   **MongoDB Cloud Account:**  Create an account at [MongoDB Cloud](https://www.mongodb.com/cloud/atlas/register) and get your connection URI.
-   **Gemini API Key:** Obtain an API key from [Google AI Studio](https://ai.google.dev/).

## Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/reddy-santhu/study-bot.git
    cd study-bot
    ```

2.  **Install Go dependencies:**

    ```bash
    go mod tidy
    go get github.com/bwmarrin/discordgo
    go get go.mongodb.org/mongo-driver/mongo
    go get github.com/joho/godotenv
    go get gopkg.in/yaml.v2
    go get github.com/go-co-op/gocron
    go get github.com/google/generative-ai-go
    go get google.golang.org/api/option
    go get github.com/ledongthuc/pdf

    ```
3. **Project Structure:**


study-bot/
│── ai/                        # AI-related code (Gemini integration)
│   └── gemini.go              # Interacts with the Gemini API
│── bot/                       # Discord bot logic
│   ├── commands/              # Command handlers
│   │   ├── ask.go             # /ask command
│   │   ├── pdf.go             # /pdf, /viewpdf, /askAi commands
│   │   ├── pomo.go            # /pomo command
│   │   ├── remind.go          # /remind command
│   │   ├── study.go           # /study commands
│   │   └── summarize.go       # /summarize command
│   └── bot.go                 # Main bot setup and event handling
│── cmd/                       # Main application entry point
│   └── main.go                # Initializes and runs the bot
│── config/                    # Configuration files
│   ├── config.yaml            # Basic configuration (bot prefix, dbname)
│   └── config.go              # Loads configuration from environment variables
│── db/                        # Database interaction
│   ├── db.go                  # Database connection and helper functions
│   └── models.go              # Data models (User, StudyLog, PDFData)
│── internal/                  # Internal application logic
│   └── study/                 # Study-related functionality
│       ├── pomodoro.go        # Pomodoro timer logic
│       └── reminders.go       # Reminder scheduling logic
│── utils/                     # Utility functions
│   └── logger.go              # Logging utilities
│── .env                       # Environment variables (KEEP THIS SECRET)
│── go.mod                     # Go module definition
│── go.sum                     # Go module checksums
│── README.md                  # Project documentation
└── .gitignore                 # Files/folders to ignore in Git


## Configuration


1.  **Create a `.env` file (in the root directory):**

    ```
    BOT_TOKEN=YOUR_DISCORD_BOT_TOKEN
    MONGODB_URI=YOUR_MONGODB_URI
    GEMINI_API_KEY=YOUR_GEMINI_API_KEY
    ```

    *   **`BOT_TOKEN`:** Your Discord bot token.
    *   **`MONGODB_URI`:** Your *full* MongoDB Atlas connection string (including username, password, and cluster address, but *without* the database name at the end).  Example: `mongodb+srv://<username>:<password>@<cluster-address>/?retryWrites=true&w=majority`
    *   **`GEMINI_API_KEY`:** Your Gemini API key.

    **Important Security Notes:**

    *   **Never commit your `.env` file to version control.**  It contains sensitive credentials.  Make sure `.env` is listed in your `.gitignore` file.
    *   The `config.yaml` file *only* contains non-sensitive configuration.

## Running the Bot

1.  **Navigate to the project directory:**

    ```bash
    cd /path/to/your/study-bot
    ```

2.  **Run the bot:**

    ```bash
    go run cmd/main.go
    ```

## Command Reference

| Command                                    | Description                                                                                           |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------- |
| `/study set <goal>`                         | Sets a new study goal.                                                                               |
| `/study list`                               | Lists all current study goals.                                                                       |
| `/study remove <goal number>`                | Removes a specific study goal (use the number from `/study list`).                                  |
| `/pomo start <work time> <break time> <task>` | Starts a Pomodoro timer (work/break times in minutes).                                             |
| `/pomo stop`                                | Stops the current Pomodoro timer.                                                                    |
| `/pomo status`                              | Shows the remaining time in the current Pomodoro session.                                            |
| `/remind <time> <task>`                     | Sets a reminder for a specific task at a given time (HH:MM or YYYY-MM-DD HH:MM).                     |
| `/ask question: <question>`                   | Asks Gemini AI a question.                                                                         |
| `/summarize text: <long text>`              | Asks Gemini AI to summarize a long text passage.                                                     |
| `/pdf`                                      | Uploads a PDF file for later interaction (content stored in MongoDB).                                |
| `/viewpdf`                                  | Displays a numbered list of PDFs uploaded by the user.                                               |
| `/askAi <pdf_number> <question>`           | Asks Gemini a question related to the content of a specific PDF (using the number from `/viewpdf`).   |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------- |

