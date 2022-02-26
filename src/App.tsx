import React, { useEffect, useState } from 'react';
import { Websocket, WebsocketBuilder, WebsocketEvents } from 'websocket-ts';
import wasm from './go/main.go'
import Link from './Link';
import Messages from './Messages';
import Input from './Input';
import './App.css';

let initiator = !!!window.location.hash
let websocket: Websocket;

export interface Message {
  direction: string;
  timestamp: string;
  text: string;
}

async function encrypt(data: string, secret: string): Promise<string> {
  const buffer = new Buffer(data);
  const encryptedData = await wasm.Encrypt(buffer, buffer.length, secret);

  return Buffer.from(encryptedData).toString('hex');
}

async function decrypt(data: string, secret: string): Promise<string> {
  const buffer = Buffer.from(data, 'hex');
  const decryptedData = await wasm.Decrypt(buffer, buffer.length, secret);

  return Buffer.from(decryptedData).toString();
}

function App() {
  const [myID, setMyID] = useState("");
  const [theirID, setTheirID] = useState("");
  const [publicKey, setPublicKey] = useState("");
  const [secretKey, setSecretKey] = useState("");
  const [state, setState] = useState("pending");
  const [message, setMessage] = useState<Message>();
  const [messages, setMessages] = useState<Message[]>([]);
  const [outgoingText, setOutgoingtext] = useState("");

  async function generateSecret(key: string) {
    const secret = await wasm.ComputeSecret(key);
    setSecretKey(secret);
    console.log('Shared secret generated', secret);
  }

  const processMessage = (i: Websocket, ev: any) => {
    const message = JSON.parse(ev.data);

    switch (message.action) {
      case 'registered':
        console.log('Registered ConnectionId with signaling server', message.payload.ConnectionId)
        setState('registered');

        if(!initiator && theirID === "") {
          const urlID = window.location.hash.substring(2);
          setTheirID(urlID);
          console.log('Got their ID from URL', urlID);
        }
        break;
      case 'exchange':
        if(initiator && theirID === "") {
          console.log('Got their ID', message.payload.id);
          setTheirID(message.payload.id);
        }

        console.log('Got their public key', message.payload.publicKey);
        generateSecret(message.payload.publicKey);

        setState('exchanged');
        break;
      case 'switch':
        console.log('Attempting to switch to WebRTC');
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

  useEffect(() => {
    if(message) {
      decrypt(message.text, secretKey).then(decrypted => {
        message.text = decrypted;
        setMessages((prevMessages: Message[]) => [...prevMessages, message]);
      }).catch(err => {
        console.log('Error decrypting message', err);
      });
    }
  }, [message, secretKey]);

  useEffect(() => {
    if(myID && theirID && publicKey && secretKey === "") {
      console.log(`Sending my publicKey to ${theirID}`);
      websocket.send(JSON.stringify({
        action: 'exchange',
        payload: {
          myID,
          theirID,
          publicKey
        }
      }));
    }
  }, [myID, theirID, publicKey, secretKey]);

  useEffect(() => {
    async function init() {
      const id = await wasm.GenerateRandomString(24)
      setMyID(id);
      console.log('Generated my ID', id);

      const key = await wasm.GenerateKey();
      setPublicKey(key);
      console.log('Generated my public key', key);

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
      websocket.addEventListener(WebsocketEvents.message, processMessage);
    }

    init();
  }, []);

  const handleSubmit = (event: React.SyntheticEvent) => {
    event.preventDefault();

    if(theirID && secretKey) {
      encrypt(outgoingText, secretKey).then(encrypted => {
        websocket.send(JSON.stringify({
          action: 'send',
          payload: {
            id: theirID,
            message: encrypted
          }
        }));

        setMessages((prevMessages: Message[]) => [...prevMessages, {
          direction: 'out',
          timestamp: new Date().toLocaleString("en-us", { hour: '2-digit', minute: '2-digit' }),
          text: outgoingText
        }]);

        setOutgoingtext("");
      }).catch(err => {
        console.log('Error encrypting message', err);
      });
    } else {
      console.log('Websocket connection not ready');
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
