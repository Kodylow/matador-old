# Renegade: A Bitcoin-Powered Passthrough Server

Renegade lets you sell API access against an arbitrary API using your API Key in exchange for bitcoin micropayments via L402 Payment Required Codes. 

I built Renegade because I'm sick of waiting for companies to wrap their APIs with Bitcoin payments, so this lets anyone with an API Key act as an L402 passthrough to the actual API, letting anyone pay for using your API Key with bitcoin.

This first version of Renegade is configured to run against the OPENAI API and currently supports the following endpoints:

```bash
POST $API_ROOT/v1/chat/completions
GET $API_ROOT/v1/images/generations
GET /v1/models
GET /v1/models/{model}
```

Support for these other endpoints will be added over the next few days:
```bash
POST /v1/images/edits
POST /v1/images/variations
POST /v1/embeddings
POST /v1/audio/transcriptions
POST /v1/audio/translations
GET /v1/files
POST /v1/files
DELETE /v1/files/{file}
GET /v1/files/{file}
GET /v1/files/content
POST /v1/fine-tunes
GET /v1/fine-tunes
GET /v1/fine-tunes/{model}
POST /v1/fine-tunes/cancel
GET /v1/fine-tunes/events
DELETE /v1/models/{model}
POST /v1/moderations
```


Renegade passes the request through exactly as if you were hitting against the actual API, replacing the L402 Authorization Header the client hits against renegade with your API key. Clients pay you in Bitcoin, you pay the API service with your credit card.

Renegade is a WIP, use at your own risk (MIT LICENSE copied below)

## Getting Started

Here's how to get Renegade up and running

### Prerequisites

You need to have Golang installed: [https://golang.org/doc/install](https://golang.org/doc/install)

Or just load this into Replit, the default configs from the checked in .replit and replit.nix work out of the box

### Clone the repository

To clone the Renegade repository to your local system, execute the following command in your terminal:

```bash
git clone https://github.com/kodylow/renegade
```

### Configuration
Post-cloning, navigate to the project root and create a .env file (or on Replit set these in Secrets). This file must include your API key, the API root, your Lightning address, and a Rune secret as follows:
```dotenv
API_KEY = YOUR_OPENAI_API_KEY
API_ROOT = "https://api.openai.com"
LN_ADDRESS = "yourusername@getalby.com"
RUNE_SECRET = "rqV9+bCcwGVNh2MkzoHnkGAp0YLrySRd1nLAlnNqrAc="
```

To generate the rune secret you just need some random base64 bytes, you can use this command: openssl rand -base64 32

You can change the pricing and endpoints as well, the current configuration is extremely conservative (will overcharge in bitcoin terms) and hardcodes a price of bitcoin at $28,000 until I get around to creating a bitcoin price service.

# Running Renegade
To launch the server, execute the following command:

```bash
go run main.go
```
Voila! Your Renegade server is now live, ready to process requests and exchange API key access for Bitcoin payments.

# MIT License

Copyright 2023 Kody Low

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
