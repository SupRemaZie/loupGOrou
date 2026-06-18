import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import CreateGame from './pages/CreateGame';
import GamePage from './pages/GamePage';
import AISpectator from './pages/AISpectator';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/create" element={<CreateGame />} />
        <Route path="/game/:id" element={<GamePage />} />
        <Route path="/ai" element={<AISpectator />} />
      </Routes>
    </BrowserRouter>
  );
}
