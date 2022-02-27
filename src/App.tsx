import React, { useEffect, useState, useRef } from 'react';
import {
  Websocket,
  WebsocketBuilder,
  WebsocketEvents,
  ConstantBackoff,
  LRUBuffer,
} from 'websocket-ts';
import SimplePeer, { SignalData } from 'simple-peer';
import {
  encrypt,
  decrypt,
  ComputeSecret,
  GenerateRandomString,
  GenerateKeyPair,
} from './lib/crypto';
import Link from './Link';
import Messages from './Messages';
import Input from './Input';
import Footer from './Footer';
import './App.css';

let initiator = !!!window.location.hash;
let websocket: Websocket;
let peer: SimplePeer.Instance;

export interface Message {
  direction: string;
  timestamp: string;
  text: string;
}

function App() {
  const [myID, setMyID] = useState<string>();
  const [theirID, setTheirID] = useState<string>();
  const [publicKey, setPublicKey] = useState<string>();
  const [secretKey, setSecretKey] = useState<string>();
  const [messages, setMessages] = useState<Message[]>([]);
  const [outgoingText, setOutgoingtext] = useState<string>('');
  const [status, setStatus] = useState<string>('Waiting for connection...');
  const [signalData, setSignalData] = useState<SignalData>();
  const [peerConnected, setPeerConnected] = useState<boolean>(false);

  function getHashData(): { theirId: string; myId: string } {
    const hash = window.location.hash;
    const keys = hash.substr(1).split('/');

    return {
      theirId: keys[1],
      myId: keys[2],
    };
  }

  function handleSubmit(event: React.SyntheticEvent) {
    event.preventDefault();

    if (theirID && secretKey && outgoingText) {
      sendMessage(outgoingText, secretKey);
    }
  }

  async function generateSecret(key: string) {
    const secret = await ComputeSecret(key);

    setSecretKey(secret);
    setStatus('Connected via server');
    console.log('Shared secret generated');
  }

  async function receiveMessage(data: string, secretKey: string) {
    const decrypted = await decrypt(data, secretKey);

    setMessages((prevMessages: Message[]) => [
      ...prevMessages,
      {
        direction: 'in',
        timestamp: new Date().toLocaleString('en-us', {
          hour: '2-digit',
          minute: '2-digit',
        }),
        text: decrypted,
      },
    ]);
  }

  async function sendMessage(data: string, secretKey: string) {
    const encrypted = await encrypt(data, secretKey);
    const msg = JSON.stringify({
      action: 'message',
      payload: {
        id: theirID,
        message: encrypted,
      },
    });

    if (peerConnected) {
      try {
        peer.send(msg);
      } catch (error) {
        console.log('Error sending via WebRTC', error);
        websocket.send(msg);
      }
    } else {
      websocket.send(msg);
    }

    setMessages((prevMessages: Message[]) => [
      ...prevMessages,
      {
        direction: 'out',
        timestamp: new Date().toLocaleString('en-us', {
          hour: '2-digit',
          minute: '2-digit',
        }),
        text: outgoingText,
      },
    ]);

    setOutgoingtext('');
  }

  useEffect(() => {
    if (signalData && secretKey) {
      console.log('Attempting to switch to WebRTC');
      encrypt(JSON.stringify(signalData), secretKey)
        .then((encrypted) => {
          websocket.send(
            JSON.stringify({
              action: 'switch',
              payload: {
                id: theirID,
                data: encrypted,
              },
            })
          );
        })
        .catch((err) => {
          console.log('Error encrypting signal data', err);
        });
    }
  }, [signalData, secretKey, theirID]);

  useEffect(() => {
    if (myID && theirID && publicKey) {
      window.location.hash = `/${theirID}/${myID}`;
      console.log(`Sending my publicKey to ${theirID}`);
      websocket.send(
        JSON.stringify({
          action: 'exchange',
          payload: {
            myID,
            theirID,
            publicKey,
          },
        })
      );
    }
  }, [myID, theirID, publicKey, secretKey]);

  const processMessage = (data: any) => {
    const message = JSON.parse(data);
    const theirUrlID = getHashData().theirId;

    switch (message.action) {
      case 'registered':
        console.log('Registered ConnectionId with signaling server', message.payload.ConnectionId);

        if (theirUrlID) {
          setTheirID(theirUrlID);
          console.log('Got their ID from URL', theirUrlID);
        }
        break;
      case 'exchange':
        if (!theirUrlID) {
          console.log('Got their ID', message.payload.id);
          setTheirID(message.payload.id);
        }

        console.log('Got their public key', message.payload.publicKey);
        generateSecret(message.payload.publicKey);
        break;
      case 'switch':
        console.log('Got request to switch to WebRTC');

        if (secretKey) {
          decrypt(message.payload.data, secretKey)
            .then((decrypted) => {
              const signalData = JSON.parse(decrypted);
              peer.signal(signalData);
            })
            .catch((err) => {
              console.log('Error decrypting signal data', err);
            });
        }
        break;
      case 'message':
        if (secretKey) {
          receiveMessage(message.payload.message, secretKey);
        }
        break;
      default:
        console.log('Unknown message action', message);
    }
  };

  const processMessageRef = useRef(processMessage);
  useEffect(() => {
    processMessageRef.current = processMessage;
  });

  useEffect(() => {
    async function init() {
      let id: string;
      const myUrlId = getHashData().myId;

      if (myUrlId) {
        id = myUrlId;
        console.log('Got my ID from URL', id);
      } else {
        id = await GenerateRandomString(24);
        console.log('Generated my ID', id);
      }

      setMyID(id);

      const publicKey = await GenerateKeyPair();

      setPublicKey(publicKey);
      console.log('Generated my public key', publicKey);

      websocket = new WebsocketBuilder('wss://api.chatsh.it')
        .withBackoff(new ConstantBackoff(1000))
        .withBuffer(new LRUBuffer(10))
        .onError((i, event) => console.log('Websocket error', event))
        .onClose(() => setStatus('Disconnected'))
        .build();

      const register = () => {
        console.log('Registering with signaling server');

        websocket.send(
          JSON.stringify({
            action: 'register',
            payload: {
              id,
            },
          })
        );
      };

      websocket.addEventListener(WebsocketEvents.open, register);
      websocket.addEventListener(WebsocketEvents.retry, register);
      websocket.addEventListener(WebsocketEvents.message, (ws: Websocket, ev: MessageEvent<any>) =>
        processMessageRef.current(ev.data)
      );

      peer = new SimplePeer({ initiator, trickle: false })
        .on('signal', (data) => {
          console.log('Got WebRTC signal data');
          setSignalData(data);
        })
        .on('error', (err) => {
          console.log('WebRTC error', err);
        })
        .on('connect', () => {
          console.log('WebRTC peer connected');
          setStatus('Connected directly with WebRTC');
          setPeerConnected(true);
        })
        .on('close', () => setPeerConnected(false))
        .on('data', (data: any) => processMessageRef.current(data));
    }

    init();
  }, []);

  return (
    <div className="body">
      <div className="App">
        <div className="container">
          <div className="row logo">
            <div className="col-md-12 text-center">
              <h1 className="h1">
                <a href="/">
                  chat<span>sh.it</span>
                </a>
              </h1>
            </div>
          </div>
          <Link id={myID} display={!secretKey && initiator} />
          <Messages messages={messages} />
          <Input
            outgoingText={outgoingText}
            setOutgoingtext={setOutgoingtext}
            handleSubmit={handleSubmit}
            display={!!secretKey}
          />
        </div>
      </div>
      <Footer status={status} />
    </div>
  );
}

export default App;
