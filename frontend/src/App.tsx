import axios from 'axios';
import './App.css'
import { useEffect, useState } from 'react';

interface Score {
  rank: number;
  score: number;
  user_id: string;
  username: string; // Added username to the Score interface
  game_id: string;
  game_name: string; // Added game_name to the Score interface
}

function App() {
  const [scores, setScores] = useState<Score[]>([]);

  useEffect(() => {
    const getScores = async () => {
      const fetchedScores = await fetchScores();
      setScores(fetchedScores);
    };
    getScores();
  }, []);

  const groupedScores = scores.reduce((acc, score) => {
    if (!acc[score.game_id]) {
      acc[score.game_id] = { game_name: score.game_name, scores: [] }; // Store game_name with scores
    }
    acc[score.game_id].scores.push(score);
    return acc;
  }, {} as Record<string, { game_name: string; scores: Score[] }>);

  return (
    <>
      <div>
        <h1>Leaderboard</h1>
        <div>
          <input type="text" placeholder="Game ID" id="gameId" />
          <input type="text" placeholder="User ID" id="userId" />
          <button onClick={handleAddScore}>Add Score</button>
        </div>
        <div>
          <h2>Scores</h2>
          {Object.keys(groupedScores).map((gameId) => (
            <div key={gameId}>
              <h3>Game: {groupedScores[gameId].game_name} (ID: {gameId})</h3> 
              <ul>
                {groupedScores[gameId].scores.map((score) => (
                  <li key={score.user_id}>
                    Rank: {score.rank} - Score: {score.score} - User: {score.username} (ID: {score.user_id})
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}

const handleAddScore = async () => {
  const gameId = (document.getElementById('gameId') as HTMLInputElement).value;
  const userId = (document.getElementById('userId') as HTMLInputElement).value;

  try {
    const response = await axios.post('http://localhost:8080/scores', {
      gameId,
      userId,
    });
    console.log('Score added successfully:', response.data);
  } catch (error) {
    console.error('Error adding score:', error);
  }
  
  (document.getElementById('gameId') as HTMLInputElement).value = '';
  (document.getElementById('userId') as HTMLInputElement).value = '';
};

const fetchScores = async (): Promise<Score[]> => {
  try {
    const response = await axios.get('http://localhost:8080/scores');
    return response.data;
  } catch (error) {
    console.error('Error fetching scores:', error);
    return [];
  }
};
export default App
