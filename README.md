# Overview
This is a 1 vs 1 turn-based minesweeper introduced more than a decade ago by MSN Messenger. The objective of the game is to reveal more mines than your opponent. Player who reveals more than half of the mines first wins the game.

# Getting Started

## Prerequisite
This is a full-stack application built with Go backend and Vue.js frontend. To run it, Go and Node.js have to first be installed.

### Go
- Follow the instructions provided by the [Go](https://go.dev/doc/install)

### Node.js
- For Windows and Mac users, download the installer from [Node.js Official Page](https://nodejs.org/en/download) 

- For Linux users, use the following link to install Node.js via package manager: [Installing Node.js via package manager](https://nodejs.org/en/download/package-manager)

## Installation
1. Clone this repository
```
git clone https://github.com/alex-fok/minesweeper
```
2. Open cloned directory, change directory to /client/minesweeper, then install npm packages
```
npm install
```


## Run in Dev

### Server
At /server, run
```
go run .
```

### Client
At /client/minesweeper, run
```
npm run dev
```
