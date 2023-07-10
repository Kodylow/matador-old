# Renegade: A Bitcoin-Powered Passthrough Server

Renegade lets you sell API access against an arbitrary API using your API Key in exchange for bitcoin micropayments via L402 Payment Required Codes. 

I built Renegade because I'm sick of waiting for companies to wrap their APIs with Bitcoin payments, so this lets anyone with an API Key act as an L402 passthrough to the actual API, letting anyone pay for using your API Key with bitcoin.

Renegade passes the request through exactly as if you were hitting against the actual API, replacing the L402 Authorization Header the client hits against renegade with your API key. Clients pay you in Bitcoin, you pay the API service with your credit card.

Renegade is a WIP, use at your own risk (MIT LICENSE copied below)

## Getting Started

Here's how to get Renegade up and running

### Prerequisites

You need to have the Go programming language installed on your system. If you haven't installed it yet, follow the official Go installation instructions here: [https://golang.org/doc/install](https://golang.org/doc/install)

Or just load this into Replit.

### Clone the repository

To clone the Renegade repository to your local system, execute the following command in your terminal:

```bash
git clone https://github.com/kodylow/renegade
```

Configuration
Post-cloning, navigate to the project root and create a .env file. This file must include your API key, Lightning Network address, and API root as follows:
```dotenv
API_KEY=your_api_key
LNADDRESS=your_lightning_network_address
API_ROOT=your_api_root
```

You'll also need to set pricing through the PricingService. You can match pricing functions for specific endpoints or use a flatfee.

Running Renegade
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
