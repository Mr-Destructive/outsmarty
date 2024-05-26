import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import webSocketService from '../websocket';
import Logout from './Logout';

const Room = () => {
  const { slug } = useParams();
  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState('');
  const [username, setUsername] = useState('');
  console.log(document.cookie);

  useEffect(() => {
    const storedUsername = document.cookie
      .split('; ')
      .find((row) => row.startsWith('outsmarty_name='))
      ?.split('=')[1];
    setUsername(storedUsername || 'Anonymous');

    const uid = document.cookie
      .split('; ')
      .find((row) => row.startsWith('outsmarty_uid='))
      ?.split('=')[1];

    webSocketService.connect(slug, uid);

    webSocketService.on('message', (payload) => {
      setMessages((prevMessages) => [...prevMessages, payload]);
    });

    return () => {
      webSocketService.disconnect();
    };
  }, [slug]);

  const handleSendMessage = () => {
    if (message.trim()) {
      const newMessage = { username, content: message, timestamp: new Date().toISOString() };
      webSocketService.send('message', newMessage);
      setMessage('');
    }
  };

  const handleKeyPress = (event) => {
    if (event.key === 'Enter') {
      handleSendMessage();
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Logout />
      <h1 className="text-4xl font-bold mb-8">Room: {slug}</h1>
      <div className="border p-4 w-3/4 h-2/3 flex flex-col">
        <div className="flex-grow overflow-y-scroll mb-4">
          {messages.map((msg, index) => (
            <div key={index} className="mb-2">
              <span className="font-bold">{msg.username}: </span>
              <span>{msg.content}</span>
            </div>
          ))}
        </div>
        <div className="flex">
          <input
            type="text"
            className="border p-2 flex-grow mr-2"
            placeholder="Type a message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyPress={handleKeyPress}
          />
          <button
            className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-700"
            onClick={handleSendMessage}
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
};

export default Room;
