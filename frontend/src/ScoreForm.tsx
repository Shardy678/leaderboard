import React, { FormEvent } from 'react';

interface Game {
    id: string;
    name: string;
}

interface ScoreFormProps {
    games: Game[];
    onSubmit: (gameId: string, userId: string, score: number) => void;
}

export const ScoreForm: React.FC<ScoreFormProps> = ({ games, onSubmit }) => {
    const handleAddScore = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const formData = new FormData(e.currentTarget);
        const gameId = formData.get('gameId') as string;
        const userId = formData.get('userId') as string;
        const score = Number(formData.get('score'));

        onSubmit(gameId, userId, score);
        e.currentTarget.reset();
    };

    return (
        <div>
            <h2>Add a new score</h2>
            <form onSubmit={handleAddScore}>
                <select name="gameId" required>
                    <option value="">Select Game</option>
                    {games.map(game => (
                        <option key={game.id} value={game.id}>{game.name}</option>
                    ))}
                </select>
                <input type="text" name="userId" placeholder="User ID" required />
                <input type="number" name="score" placeholder="Score" required />
                <button type="submit">Add Score</button>
            </form>
        </div>
    );
};