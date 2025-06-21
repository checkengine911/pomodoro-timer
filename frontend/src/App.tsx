import React, { useState, useRef, useEffect } from 'react';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || '';

function Login({ onLogin, onShowRegister }: { onLogin: () => void, onShowRegister: () => void }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || 'Ошибка авторизации');
        return;
      }
      const data = await res.json();
      localStorage.setItem('jwt', data.token);
      onLogin();
    } catch (e) {
      setError('Ошибка сети');
    }
  };

  return (
    <div className="login-container">
      <h2>Вход</h2>
      <form onSubmit={handleSubmit}>
        <input type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} required />
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} required />
        <button type="submit">Войти</button>
        {error && <div className="error">{error}</div>}
      </form>
      <div className="switch-auth">
        Нет аккаунта? <button type="button" onClick={onShowRegister} className="link-btn">Зарегистрироваться</button>
      </div>
    </div>
  );
}

function Register({ onRegister }: { onRegister: () => void }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (password !== confirm) {
      setError('Пароли не совпадают');
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || 'Ошибка регистрации');
        setLoading(false);
        return;
      }
      // После успешной регистрации — сразу логиним пользователя
      const loginRes = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      if (loginRes.ok) {
        const data = await loginRes.json();
        localStorage.setItem('jwt', data.token);
        onRegister();
      } else {
        setError('Регистрация прошла, но не удалось войти');
      }
    } catch (e) {
      setError('Ошибка сети');
    }
    setLoading(false);
  };

  return (
    <div className="login-container">
      <h2>Регистрация</h2>
      <form onSubmit={handleSubmit}>
        <input type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} required />
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} required />
        <input type="password" placeholder="Повторите пароль" value={confirm} onChange={e => setConfirm(e.target.value)} required />
        <button type="submit" disabled={loading}>{loading ? 'Регистрация...' : 'Зарегистрироваться'}</button>
        {error && <div className="error">{error}</div>}
      </form>
      <div className="switch-auth">
        Уже есть аккаунт? <button type="button" onClick={onRegister} className="link-btn">Войти</button>
      </div>
    </div>
  );
}

const POMODORO_DURATION = 25 * 60; // 25 минут
const SHORT_BREAK = 5 * 60; // 5 минут
const LONG_BREAK = 15 * 60; // 15 минут
const CYCLES_BEFORE_LONG_BREAK = 4;

type TimerMode = 'pomodoro' | 'short_break' | 'long_break';

const modeLabels: Record<TimerMode, string> = {
  pomodoro: 'Помодоро',
  short_break: 'Короткий перерыв',
  long_break: 'Длинный перерыв',
};

type PomodoroSession = {
  id: number;
  duration: number;
  start_time: string;
  end_time: string;
  task_id?: number;
};

function Pomodoro() {
  const [secondsLeft, setSecondsLeft] = useState(POMODORO_DURATION);
  const [isRunning, setIsRunning] = useState(false);
  const [mode, setMode] = useState<TimerMode>('pomodoro');
  const [cycle, setCycle] = useState(0);
  const [sessions, setSessions] = useState<PomodoroSession[]>([]);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  // Получение сессий с бекенда
  useEffect(() => {
    const fetchSessions = async () => {
      const jwt = localStorage.getItem('jwt');
      if (!jwt) return;
      const res = await fetch(`${API_URL}/pomodoro`, {
        headers: { 'Authorization': `Bearer ${jwt}` }
      });
      if (res.ok) {
        const data = await res.json();
        setSessions(data);
      }
    };
    fetchSessions();
  }, []);

  useEffect(() => {
    if (isRunning) {
      intervalRef.current = setInterval(() => {
        setSecondsLeft((prev) => {
          if (prev > 0) return prev - 1;
          handleTimerEnd();
          return 0;
        });
      }, 1000);
    } else if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
    // eslint-disable-next-line
  }, [isRunning, mode]);

  const handleTimerEnd = async () => {
    setIsRunning(false);
    const jwt = localStorage.getItem('jwt');
    const now = new Date();
    let duration = 0;
    if (mode === 'pomodoro') duration = POMODORO_DURATION;
    if (mode === 'short_break') duration = SHORT_BREAK;
    if (mode === 'long_break') duration = LONG_BREAK;
    // Создаём сессию на бекенде
    if (jwt && mode === 'pomodoro') {
      await fetch(`${API_URL}/pomodoro`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${jwt}`
        },
        body: JSON.stringify({
          duration: duration / 60,
          start_time: new Date(now.getTime() - duration * 1000).toISOString(),
          end_time: now.toISOString(),
        })
      });
    }
    if (mode === 'pomodoro') {
      const nextCycle = cycle + 1;
      setCycle(nextCycle);
      if (nextCycle % CYCLES_BEFORE_LONG_BREAK === 0) {
        setMode('long_break');
        setSecondsLeft(LONG_BREAK);
      } else {
        setMode('short_break');
        setSecondsLeft(SHORT_BREAK);
      }
    } else {
      setMode('pomodoro');
      setSecondsLeft(POMODORO_DURATION);
    }
  };

  const handleStart = () => {
    setIsRunning(true);
  };

  const handlePause = () => {
    setIsRunning(false);
  };

  const handleReset = () => {
    setIsRunning(false);
    if (mode === 'pomodoro') setSecondsLeft(POMODORO_DURATION);
    if (mode === 'short_break') setSecondsLeft(SHORT_BREAK);
    if (mode === 'long_break') setSecondsLeft(LONG_BREAK);
  };

  const minutes = String(Math.floor(secondsLeft / 60)).padStart(2, '0');
  const seconds = String(secondsLeft % 60).padStart(2, '0');

  return (
    <div className="pomodoro-container">
      <h1>Pomodoro Таймер</h1>
      <div className="mode-switcher">
        {(['pomodoro', 'short_break', 'long_break'] as TimerMode[]).map((m) => (
          <button
            key={m}
            className={mode === m ? 'active' : ''}
            onClick={() => {
              setMode(m);
              setIsRunning(false);
              if (m === 'pomodoro') setSecondsLeft(POMODORO_DURATION);
              if (m === 'short_break') setSecondsLeft(SHORT_BREAK);
              if (m === 'long_break') setSecondsLeft(LONG_BREAK);
            }}
          >
            {modeLabels[m]}
          </button>
        ))}
      </div>
      <div className="timer-display">
        <span>{minutes}:{seconds}</span>
      </div>
      <div className="controls">
        {!isRunning ? (
          <button onClick={handleStart} className="start">Старт</button>
        ) : (
          <button onClick={handlePause} className="pause">Пауза</button>
        )}
        <button onClick={handleReset} className="reset">Сброс</button>
      </div>
      <div className="cycle-info">
        Цикл: {cycle % CYCLES_BEFORE_LONG_BREAK} / {CYCLES_BEFORE_LONG_BREAK}
      </div>
      <div className="sessions-list">
        <h3>Ваши сессии</h3>
        <ul>
          {sessions.map(s => (
            <li key={s.id}>
              {new Date(s.start_time).toLocaleString()} — {s.duration} мин
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

function App() {
  const [isAuth, setIsAuth] = useState(!!localStorage.getItem('jwt'));
  const [showRegister, setShowRegister] = useState(false);

  if (isAuth) return <Pomodoro />;
  if (showRegister) return <Register onRegister={() => { setIsAuth(true); setShowRegister(false); }} />;
  return <Login onLogin={() => setIsAuth(true)} onShowRegister={() => setShowRegister(true)} />;
}

export default App;
