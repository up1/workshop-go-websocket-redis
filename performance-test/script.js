import ws from 'k6/ws';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 30 },
    { duration: '30s', target: 30 },
    { duration: '20s', target: 60 },
    { duration: '30s', target: 60 },
    { duration: '20s', target: 30 },
    { duration: '30s', target: 30 },
    { duration: '10s', target: 0 },
  ],
};

export default function () {
  const text = 'Hello from server, Guest!';
  const url = 'ws://localhost:8080/ws/user01';

  const res = ws.connect(url, null, function (socket) {
    socket.on('open', function open() {
      console.log('connected');
      socket.setInterval(function interval() {
        socket.send(text);
        console.log('Message sent: ', text);
      }, 1000);
    });

    socket.on('message', function message(data) {
      console.log('Message received: ', data);
      check(data, { 'data is correct': (r) => r && r === text });
    });

    socket.on('close', () => console.log('disconnected'));

    socket.setTimeout(function () {
      console.log('5 seconds passed, closing the socket');
      socket.close();
    }, 5000);
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}