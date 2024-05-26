import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createRoom } from '../api'; // Import your API function

const CreateRoom = () => {
  const [roomName, setRoomName] = useState('');
  const navigate = useNavigate();

  const handleCreateRoom = async () => {
    try {
      // Call your API function to create the room
      const response = await createRoom(roomName);
      
      // Extract the room slug from the response or use the room name
      const roomSlug = response.roomSlug;

      // Redirect to the room component
      navigate(`/room/${roomSlug}`);
    } catch (error) {
      console.error('Error creating room', error);
    }
  };

  return (
    <div>
      {/* Room creation form */}
      <input
        type="text"
        placeholder="Room Name"
        value={roomName}
        onChange={(e) => setRoomName(e.target.value)}
      />
      <button onClick={handleCreateRoom}>Create Room</button>
    </div>
  );
};

export default CreateRoom;
