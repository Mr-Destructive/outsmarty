import React from 'react';
import { useNavigate } from 'react-router-dom';

const Logout = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    document.cookie = 'outsmarty_uid=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';

    navigate('/');
  };

  return (
      <button
        className="bg-red-500 text-white py-2 px-4 rounded hover:bg-red-700"
        onClick={handleLogout}
      >
        Logout
      </button>
  );
};

export default Logout;

