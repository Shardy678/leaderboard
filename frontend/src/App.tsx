import axios from 'axios';
import './App.css';
import { useEffect, useState } from 'react';

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
  const [scores, setScores] = useState<Set<Score>>(new Set());

  useEffect(() => {
    const fetchGames = async () => {
        try {
            const response = await axios.get('http://localhost:8080/games');
            setGames(response.data);
        } catch (error) {
            console.error('Error fetching games:', error);
        }
    };

    fetchGames();
  }, []); 

  const fetchScores = async () => {
    const allScores: Score[] = [];
    for (const game of games) {
        try {
            const response = await axios.get(`http://localhost:8080/scores/${game.id}`);
            response.data.forEach((scoreData: { member: string; score: number; score_id: string; }) => {
                allScores.push({ 
                    score_id: scoreData.score_id,
                    user_id: scoreData.member,
                    score: scoreData.score
                });
            });
        } catch (error) {
            console.error(`Error fetching scores for game ${game.id}:`, error);
        }
    }
    setScores(new Set(allScores));
  };

  useEffect(() => {
    if (games.length > 0) {
        fetchScores();
    }
  }, [games]); 

  const handleAddScore = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.target as HTMLFormElement);
    const gameId = formData.get('gameId') as string;
    const userId = formData.get('userId') as string;
    const score = parseInt(formData.get('score') as string);
    try {
      const response = await axios.post("http://localhost:8080/scores", {
        game_id: gameId,
        user_id: userId,
        score: score
      });
      await fetchScores(); 
    } catch (error) {
      console.error('Error adding score:', error);
    }
  };

  const handleDeleteScore = async (scoreId: string, userId: string) => {
    try {
      await axios.delete(`http://localhost:8080/scores/${scoreId}/${userId}`); 
      await fetchScores(); 
    } catch (error) {
      console.error('Error deleting score:', error);
    }
  };

  return (
    <>
      <div>
        <h1>Leaderboard</h1>
        <div>
          <h2>Add a new score</h2>
          <form onSubmit={handleAddScore}>
            <input type="text" name="gameId" placeholder="Game ID" />
            <input type="text" name="userId" placeholder="User ID" />
            <input type="number" name="score" placeholder="Score" />
            <button type="submit">Add Score</button>
          </form>
        </div>
        <div>
          <h2>Scores</h2>
          {games.map(game => {
            const gameScores = Array.from(scores).filter(score => score.score_id === game.id);
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
