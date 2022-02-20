import React, { useEffect, useState } from 'react';
import { Websocket, WebsocketBuilder, WebsocketEvents } from 'websocket-ts';
import Peer from 'simple-peer';
import wasm from './go/main.go'
import Link from './Link';
import Messages from './Messages';
import Input from './Input';
import './App.css';

let theirID: string;
let secret: string;
let peer: Peer.Instance;
let initiator = !!!window.location.hash
let websocket: Websocket;

export interface Message {
  direction: string;
  timestamp: string;
  text: string;
}

async function encrypt(data: string): Promise<string> {
  const buffer = new Buffer(data);
  const encryptedData = await wasm.Encrypt(buffer, buffer.length, secret);

  return encryptedData;
}

async function createPeer() {
  peer = new Peer({
    initiator: initiator,
    trickle: false
  });

  peer.on('error', err => console.log('peer error', err))

  peer.on('signal', data => {
    encrypt(JSON.stringify(data));
    //setMessages((prevMessages: string[]) => [...prevMessages, "Got signal data"])
  })

  peer.on('connect', async () => {
    console.log('WebRTC connected');
  })

  peer.on('data', async (data) => {
    const string = new TextDecoder("utf-8").decode(data);

    console.log('WebRTC data', string);
  })
}

function App() {
  const [myID, setMyID] = useState("");
  const [publicKey, setPublicKey] = useState("");
  const [state, setState] = useState("pending");
  const [messages, setMessages] = useState<Message[]>([]);
  const [outgoingText, setOutgoingtext] = useState("");

  function setMessage(message: Message) {
    setMessages((prevMessages: Message[]) => [...prevMessages, message])
  }

  async function generateSecret(key: string) {
    secret = await wasm.ComputeSecret(key);
    console.log('Secret generated', secret)
  }

  useEffect(() => {
    const processMessage = (i: Websocket, ev: any) => {
      const message = JSON.parse(ev.data);

      switch (message.action) {
        case 'registered':
          console.log('Registered with signaling server', message)
          setState('registered');
          if(!initiator) {
            theirID = window.location.hash.substring(2);
            console.log(`Sending my publicKey to ${theirID}`);
            i.send(JSON.stringify({
              action: 'exchange',
              payload: {
                myID,
                theirID,
                publicKey
              }
            }));
          }
          break;
        case 'exchange':
          console.log('Received key', message)

          if(!theirID) {
            theirID = message.payload.id;
            console.log(`Sending my publicKey back to ${theirID}`);
            i.send(JSON.stringify({
              action: 'exchange',
              payload: {
                myID,
                theirID,
                publicKey
              }
            }));
          }

          generateSecret(message.payload.publicKey);
          setState('exchanged');
          break;
        case 'switch':
          console.log('Attempting to switch to WebRTC')
          peer.signal(message.payload);
          break;
        case 'message':
          setMessage({
            direction: 'in',
            timestamp: new Date().toLocaleString("en-us", { hour: '2-digit', minute: '2-digit' }),
            text: message.payload.message
          })
          break;
        default:
          console.log('Unknown message action', message)
      }
    };

    if(myID !== "" && publicKey !== "") {
      websocket.addEventListener(WebsocketEvents.message, processMessage);
    }
  }, [myID, publicKey]);

  useEffect(() => {
    async function init() {
      const id = await wasm.GenerateRandomString(24)
      setMyID(id);

      const key = await wasm.GenerateKey();
      setPublicKey(key);

      console.log('Got ID', id);

      websocket = new WebsocketBuilder('wss://api.chatsh.it')
        .onError((i, ev) => console.log('Websocket error', ev))
        .build();

      const register = () => {
        console.log('Registering with signaling server')

        websocket.send(JSON.stringify({
          action: 'register',
          payload: {
            id,
          }
        }));
      }

      websocket.addEventListener(WebsocketEvents.open, register);
      websocket.addEventListener(WebsocketEvents.retry, register);
    }

    init();
  }, []);

  const handleSubmit = (event: React.SyntheticEvent) => {
    event.preventDefault();

    if(theirID) {
      encrypt(outgoingText).then(encrypted => {
        websocket.send(JSON.stringify({
          action: 'send',
          payload: {
            id: theirID,
            message: encrypted
          }
        }));

        setMessage({
          direction: 'out',
          timestamp: new Date().toLocaleString("en-us", { hour: '2-digit', minute: '2-digit' }),
          text: outgoingText
        });

        setOutgoingtext("");
      });
    } else {
      console.log('Connection not ready');
    }
  }

  return (
    <div className="body">
      <div className="App">
        <div className="container">
          <div className="row logo">
            <div className="col-md-12 text-center">
              <h1 className="h1"><a href="/">chat<span>sh.it</span></a></h1>
            </div>
          </div>
          <Link id={ myID } display={ state === "registered" && initiator } />
          <Messages messages={ messages } />
          <Input outgoingText={ outgoingText } setOutgoingtext={ setOutgoingtext } handleSubmit={ handleSubmit } display={ state === "exchanged" } />
        </div>
      </div>
    </div>
  );
}

export default App;
