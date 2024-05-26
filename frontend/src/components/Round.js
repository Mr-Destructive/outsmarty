
import React, { useState } from 'react';

const Round = ({ nextRound }) => {
  const [answer, setAnswer] = useState('');
  
  const submitAnswer = () => {
    // Logic to handle answer submission
    console.log("Answer submitted:", answer);
    nextRound();
  };

  return (
    <div className="flex flex-col items-center">
      <p className="mb-4">Who is the fastest land animal?</p>
      <input
        type="text"
        className="border p-2 mb-4"
        value={answer}
        onChange={(e) => setAnswer(e.target.value)}
      />
      <button
        className="bg-green-500 text-white py-2 px-4 rounded hover:bg-green-700"
        onClick={submitAnswer}
      >
        Submit Answer
      </button>
    </div>
  );
};

export default Round;
