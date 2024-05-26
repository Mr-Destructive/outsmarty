import React, { useState } from 'react';
import Round from './Round';

const Game = () => {
  const [currentRound, setCurrentRound] = useState(1);

  const nextRound = () => {
    setCurrentRound(currentRound + 1);
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-4xl font-bold mb-8">Round {currentRound}</h1>
      <Round nextRound={nextRound} />
    </div>
  );
};

export default Game;
