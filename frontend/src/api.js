const API_URL = 'http://localhost:8080';

export const register = async (name, password) => {
  const response = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, password }),
  });
  return response.json();
};

export const login = async (name, password) => {
  try {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, password }),
      credentials: 'include',
    });
    
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  } catch (error) {
    console.error('Login failed', error);
    throw error;
  }
};

export const createRoom = async (name, max_players, game_rounds) => {
  const response = await fetch(`${API_URL}/rooms/create`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, max_players, game_rounds }),
  });
  return response.json();
};

export const joinRoom = async (roomId) => {
  const response = await fetch(`${API_URL}/rooms/join`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ roomId }),
  });
  return response.json();
};

export const fetchThemes = async () => {
  const response = await fetch(`${API_URL}/theme/list`)
    console.log(response);
  return response.json();
}

export const createGame = async (theme, rounds) => {
  const response = await fetch(`${API_URL}/games/create`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
      body: JSON.stringify({ "theme_id": theme, "num_rounds": rounds }),
  });
  return response.json();
}

export const getGame = async (gameId) => {
  const response = await fetch(`${API_URL}/games/detail?game_id=${gameId}`)
  return response.json();
}

export const getQuestionsForGame = async (theme, rounds) => {
  const response = await fetch(`${API_URL}/questions/generate?rounds=${rounds}&theme=${theme}`)
  return response.json();
}
