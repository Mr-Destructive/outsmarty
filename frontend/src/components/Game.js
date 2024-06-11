import React, { useState, useEffect } from 'react';
import { getGame, getQuestionsForGame } from '../api';

const Game = ({ gameId }) => {
  const [roundNum, setRoundNum] = useState(1);
  const [totalRounds, setTotalRounds] = useState(0);
  const [theme, setTheme] = useState('');
  const [questions, setQuestions] = useState([]);
  const [answer, setAnswer] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGameDetail = async (gameId) => {
      try {
        const data = await getGame(gameId);
        setRoundNum(data.round_num);
        setTheme(data.theme.name);
        setTotalRounds(data.total_rounds);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching game details:', error);
      }
    };

    fetchGameDetail(gameId);
  }, [gameId]);

  useEffect(() => {
    const fetchQuestionsForRound = async () => {
      try {
        const data = await getQuestionsForGame(theme, roundNum);
        setQuestions(data.questions);
      } catch (error) {
        console.error('Error fetching questions for round:', error);
      }
    };

    if (loading) {
      fetchQuestionsForRound();
    }
  }, [roundNum, loading, theme]);

  const handleSubmit = () => {
    // Send answer to backend
    setRoundNum(roundNum + 1);
    setAnswer('');
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="fixed top-0 left-0 w-full h-full bg-gray-800 bg-opacity-75 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg">
        <h1 className="text-3xl font-bold mb-4">Round {roundNum} of {totalRounds}</h1>
        {questions.length > 0 && (
          <div>
            <h2 className="text-2xl font-bold mb-4">Question 1:</h2>
            <p className="mb-4">{questions[0]}</p>
          </div>
        )}
        <input
          type="text"
          placeholder="Your Answer"
          className="border p-2 w-full mb-4"
          value={answer}
          onChange={(e) => setAnswer(e.target.value)}
        />
        <button
          className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-700"
          onClick={handleSubmit}
        >
          Submit
        </button>
      </div>
    </div>
  );
};

export default Game;
