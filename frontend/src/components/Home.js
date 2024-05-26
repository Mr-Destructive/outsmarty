import React from 'react';
import { Link } from 'react-router-dom';
import { isLoggedIn } from '../auth';

const Home = () => {
  if (isLoggedIn()) {
    window.location.href = '/dashboard';
  }
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-4xl font-bold mb-8">Outsmarty</h1>
      <div className="space-x-4">
        <Link to="/register">
          <button className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-700">
            Register
          </button>
        </Link>
        <Link to="/login">
          <button className="bg-green-500 text-white py-2 px-4 rounded hover:bg-green-700">
            Login
          </button>
        </Link>
      </div>
    </div>
  );
};

export default Home;
