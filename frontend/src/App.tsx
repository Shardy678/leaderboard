import axios from 'axios';
import './App.css';
import { useEffect, useState } from 'react';
import { ScoreForm } from './ScoreForm';

interface Score {
  score_id: string;
  user_id: string;
  score: number;
}

interface Game {
  id: string;
  name: string;
}

function App() {
  const [games, setGames] = useState<Game[]>([]);
  const [scores, setScores] = useState<Score[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchGames = async () => {
    try {
      const response = await axios.get('http://localhost:8080/games');
      setGames(response.data);
    } catch (error) {
      setError('Failed to fetch games');
      console.error('Error fetching games:', error);
    }
  };

  const fetchScores = async () => {
    try {
      const response = await axios.get('http://localhost:8080/scores');
      setScores(response.data.map((scoreData: { member: string; score: number; score_id: string; }) => ({
        score_id: scoreData.score_id,
        user_id: scoreData.member,
        score: scoreData.score
      })));
    } catch (error) {
      setError('Failed to fetch scores');
      console.error('Error fetching scores:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      await fetchGames();
      await fetchScores();
    };
    loadData();
  }, []);

  const handleAddScore = async (gameId: string, userId: string, score: number) => {
    setError(null);
    try {
      await axios.post("http://localhost:8080/scores", {
        game_id: gameId,
        user_id: userId,
        score: score
      });
      await fetchScores();
    } catch (error) {
      setError('Failed to add score');
      console.error('Error adding score:', error);
    }
  };

  const handleDeleteScore = async (scoreId: string, userId: string) => {
    try {
      await axios.delete(`http://localhost:8080/scores/${scoreId}/${userId}`);
      await fetchScores();
    } catch (error) {
      setError('Failed to delete score');
      console.error('Error deleting score:', error);
    }
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <>
      <div>
        <h1>Leaderboard</h1>
        {error && <div className="error-message">{error}</div>}
        <ScoreForm games={games} onSubmit={handleAddScore}/>
        <div>
          <h2>Scores</h2>
          {games.map(game => {
            const gameScores = scores.filter(score => score.score_id === game.id);
            return (
              <div key={game.id}>
                <h3>Game: {game.name.charAt(0).toUpperCase() + game.name.slice(1)} ID: {game.id}</h3>
                {gameScores.length > 0 ? (
                  <table>
                    <thead>
                      <tr>
                        <th>User</th>
                        <th>Score</th>
                        <th>Delete</th>
                      </tr>
                    </thead>
                    <tbody>
                      {gameScores.map(score => (
                        <tr key={score.user_id}>
                          <td>{score.user_id}</td>
                          <td>{score.score}</td>
                          <td>
                            <button onClick={() => handleDeleteScore(score.score_id, score.user_id)}>❌</button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                ) : (
                  <p>No scores available for this game.</p>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </>
  );
}

export default App;
