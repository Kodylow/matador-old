# Matador: A Bitcoin-Powered Passthrough Server

Matador lets you sell API access against an arbitrary API using your API Key in exchange for bitcoin micropayments via L402 Payment Required Codes. 

I built Matador because I'm sick of waiting for companies to wrap their APIs with Bitcoin payments, so this lets anyone with an API Key act as an L402 passthrough to the actual API, letting anyone pay for using your API Key with bitcoin.

This first version of Matador is configured to run against the OPENAI API and currently supports the following endpoints:

```bash
POST $API_ROOT/v1/chat/completions
GET $API_ROOT/v1/images/generations
GET /v1/models
GET /v1/models/{model}
```

You can try it out by hitting exactly like you would hit against `https://api.openai.com` but without the OpenAI Authentication Header:

```bash
curl -k -v https://matador.kody.repl.co/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": "Hello!"}]
  }'
```

This will return an L402 token and invoice, the invoice is quoted for the previous request's parameters (method, path, body).

```bash
Www-Authenticate: L402 token=48IUkiWUzeeHsmV-fIhHBdRoeMDVEfc5WLFhYRRE_zJwYXltZW50SGFzaD04MTk1Y2YxOWJkNmQ0YTIxZTY5ZTJjYThhMmE4YTIyZGY3NjdiYTVmMzc0MmVkNmE5Njk5OTI0NWZiYTIyZjcxJnJlcXVlc3RIYXNoPWFlN2Q3ZTU0MzIzNTgzNzRmODZmNjAxZmYzYzljOTFlZTRlMWZjYjAyZTViNmU5OThkMmU1OWUzMzYzYzIwYmE, invoice=lnbc80n1pj2udc5pp5sx2u7xdad49zre579j52929z9hmk0wjlxapw665knxfyt7az9acshp56kymqtxr5es99pd82vjjnmssr2l72l379pv87d05c5pd4s2n0ysqcqzzsxqyz5vqsp50per6u35xrl3uh0ak7q0qql3mvr0ep2kr04p7d4mkgjdfnv9cw6q9qyyssqp7pvnssphg9dgh35l35jlwtpcy7lvleuqjv4u7jmczu4umnc9mukcxdq9p0n3eg4a2ezfqlux7kc47qkdp9q30cvdrkcgunr4pcnlusqh8m5e0
```

Pay the lightning invoice to get the preimage and add retry the same request with the L402 authorization header:

```bash
curl -k -v http://localhost:8080/v1/chat/completions   -H "Content-Type: application/json" -H "Authorization: L402 48IUkiWUzeeHsmV-fIhHBdRoeMDVEfc5WLFhYRRE_zJwYXltZW50SGFzaD04MTk1Y2YxOWJkNmQ0YTIxZTY5ZTJjYThhMmE4YTIyZGY3NjdiYTVmMzc0MmVkNmE5Njk5OTI0NWZiYTIyZjcxJnJlcXVlc3RIYXNoPWFlN2Q3ZTU0MzIzNTgzNzRmODZmNjAxZmYzYzljOTFlZTRlMWZjYjAyZTViNmU5OThkMmU1OWUzMzYzYzIwYmE:7660c22f7e59fba0bfce676f666bc0bb81286e8594028c7d4f8715b7d8e48297"  -d '{                            "model": "gpt-3.5-turbo",
    "messages": [{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": "Hello!"}]
  }'
```

And you'll get the standard API response from the service you're hitting against.

Olé!! You just paid bitcoin to hit the OpenAI API. Now it's actually open to all!

Support for the other endpoints for audio, embeddings, and images will be added over the next few days.

Matador passes the request through exactly as if you were hitting against the actual API, replacing the L402 Authorization Header the client hits against matador with your API key. Clients pay you in Bitcoin, you pay the API service with your credit card.

Matador is a WIP, use at your own risk (MIT LICENSE copied below)

## Getting Started

Here's how to get Matador up and running

### Prerequisites

You need to have Golang installed: [https://golang.org/doc/install](https://golang.org/doc/install)

Or just load this into Replit, the default configs from the checked in .replit and replit.nix work out of the box

### Clone the repository

To clone the Matador repository to your local system, execute the following command in your terminal:

```bash
git clone https://github.com/kodylow/matador
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

# Running Matador
To launch the server, execute the following command:

```bash
go run main.go
```
Voila! Your Matador server is now live, ready to process requests and exchange API key access for Bitcoin payments.

# MIT License

Copyright 2023 Kody Low

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
