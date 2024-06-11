import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import webSocketService from '../websocket';
import Logout from './Logout';
import Game from './Game';
import { fetchThemes, createGame } from '../api';

const Room = () => {
  const { slug } = useParams();
  const [showGameOverlay, setShowGameOverlay] = useState(false);
  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState('');
  const [username, setUsername] = useState('');
  const [gameStarted, setGameStarted] = useState(false);
  const [gameId, setGameId] = useState(null);
  const [selectedTheme, setSelectedTheme] = useState({ id: '', name: '' });
  const [rounds, setRounds] = useState(3);
  const [availableThemes, setAvailableThemes] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const storedUsername = document.cookie
      .split('; ')
      .find((row) => row.startsWith('outsmarty_name='))
      ?.split('=')[1];
    setUsername(storedUsername || '');
    if (!storedUsername) {
      navigate('/login');
    }

    const uid = document.cookie
      .split('; ')
      .find((row) => row.startsWith('outsmarty_uid='))
      ?.split('=')[1];

    webSocketService.connect(slug, uid);

    webSocketService.on('message', (payload) => {
      setMessages((prevMessages) => [...prevMessages, payload]);
    });

    fetchThemesData();
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

  const fetchThemesData = async () => {
    try {
      const data = await fetchThemes();
      if (Array.isArray(data)) {
        setAvailableThemes(data);
        if (data.length > 0) {
          setSelectedTheme({ id: data[0].id, name: data[0].name });
        }
      }
    } catch (error) {
      console.error('Error fetching themes:', error);
    }
  };

  const startGame = async (theme, rounds) => {
    try {
      const data = await createGame(theme, rounds);
      setGameId(data.id);
      setShowGameOverlay(true);
    } catch (error) {
      console.error('Error starting game:', error);
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
        <div className="flex flex-col mt-4">
          <label className="mb-2">Select Theme:</label>
          <select
            className="border p-2 mb-4"
            value={selectedTheme.id}
            onChange={(e) => {
              const selectedOption = availableThemes.find(theme => theme.id === parseInt(e.target.value));
              setSelectedTheme(selectedOption);
            }}
          >
            <option value="" disabled>Select a theme</option>
            {availableThemes.map((themeOption) => (
              <option key={themeOption.id} value={themeOption.id}>{themeOption.name}</option>
            ))}
          </select>
          <label className="mb-2">Number of Rounds:</label>
          <input
            type="number"
            className="border p-2 mb-4"
            value={rounds}
            onChange={(e) => setRounds(Number(e.target.value))}
            min="1"
            max="10"
          />
          <button
            className="bg-green-500 text-white py-2 px-4 rounded hover:bg-green-700"
            onClick={() => startGame(selectedTheme.id, rounds)}
          >
            Start Game
          </button>
        </div>
      </div>
      {showGameOverlay && <Game gameId={gameId} />}
    </div>
  );
};

export default Room;
