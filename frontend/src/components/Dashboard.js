import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { createRoom, joinRoom } from '../api';
import webSocketService from '../websocket';
import Logout from './Logout';

const Dashboard = () => {
  const [roomId, setRoomId] = useState('');
  const [roomName, setRoomName] = useState('');
  const [maxPlayers, setMaxPlayers] = useState(4);
  const [gameRounds, setGameRounds] = useState(3);
  const [uid, setUid] = useState('');
const navigate = useNavigate();

  useEffect(() => {
    const storedUid = document.cookie
      .split('; ')
      .find((row) => row.startsWith('outsmarty_uid='))
      ?.split('=')[1];
    setUid(storedUid);
  }, []);

  const handleCreateRoom = async () => {
    try {
      const response = await createRoom(roomName, maxPlayers, gameRounds);
        // get uid from cookie
        let user_id = document.cookie
        .split('; ')
        .find((row) => row.startsWith('outsmarty_uid='))
        ?.split('=')[1];
        setRoomId(response.roomId);
    } catch (error) {
      console.error('Room creation failed', error);
    }
  };

  const handleJoinRoom = async () => {
    try {
        navigate(`/room/${roomId}`);
        
    } catch (error) {
      console.error('Joining room failed', error);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Logout />
      <h1 className="text-4xl font-bold mb-8">Dashboard</h1>
      <div className="space-y-4">
        <input
          type="text"
          placeholder="Room Name"
          className="border p-2"
          value={roomName}
          onChange={(e) => setRoomName(e.target.value)}
        />
        <input
          type="number"
          placeholder="Max Players"
          className="border p-2"
          value={maxPlayers}
          onChange={(e) => setMaxPlayers(e.target.value)}
        />
        <input
          type="number"
          placeholder="Game Rounds"
          className="border p-2"
          value={gameRounds}
          onChange={(e) => setGameRounds(e.target.value)}
        />
        <button
          className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-700"
          onClick={handleCreateRoom}
        >
          Create Room
        </button>
        <input
          type="text"
          placeholder="Enter Room ID"
          className="border p-2"
          value={roomId}
          onChange={(e) => setRoomId(e.target.value)}
        />
        <button
          className="bg-green-500 text-white py-2 px-4 rounded hover:bg-green-700"
          onClick={handleJoinRoom}
        >
          Join Room
        </button>
      </div>
    </div>
  );
};

export default Dashboard;
