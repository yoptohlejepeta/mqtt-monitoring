import * as WebSocket from 'ws';

const ws = new WebSocket('ws://localhost:8000');

ws.on('open', () => {
  console.log('Connected to the server');
  ws.send('Hello from client');
});

ws.on('message', (message) => {
  console.log(`Received message: ${message}`);
});

ws.on('error', (error) => {
  console.log(`Error occurred: ${error}`);
});

ws.on('close', () => {
  console.log('Disconnected from the server');
});
