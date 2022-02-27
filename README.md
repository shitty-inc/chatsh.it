# Chatsh.it

![GitHub release](https://img.shields.io/github/release/shitty-inc/chatsh.it.svg)

> Chatsh.it is a simple peer to peer web chat app with the aim of being overly end to end encrypted

## Description

First we establish a websocket connection to our server and generate you a public key. This is in the URL you can share with your friend. Once they click on the we perform our own ECDH key-exchange using Curve25519 to generate a shared private key to use to encrypt all messages with the [TripleSec](https://github.com/keybase/go-triplesec) library. We then use the websocket connection to attempt to establish a peer to peer connection secured by DTLS using WebRTC.

## Disclaimer

I am not an expert in cryptography. If you have something important to keep secret please think about using a peer reviewed and audited service. This is just an experiment with WebRTC and WASM.

If you are behind some kinds of NAT configuration then chatsh.it WebRTC part might not work.
