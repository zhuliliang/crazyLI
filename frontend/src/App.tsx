import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Home } from './pages/Home';
import { LoginSuccess } from './pages/LoginSuccess';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login/success" element={<LoginSuccess />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
