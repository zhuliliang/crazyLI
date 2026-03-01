import { FormEvent, useEffect, useState } from 'react';
import {
  API_BASE_URL,
  fetchCurrentUser,
  loginWithEmail,
  registerWithEmail,
  type User,
} from '../lib/api';
import { UserCard } from '../components/UserCard';

export function Home() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [form, setForm] = useState({ name: '', email: '', password: '' });
  const [feedback, setFeedback] = useState<string | null>(null);

  useEffect(() => {
    fetchCurrentUser()
      .then(setUser)
      .finally(() => setLoading(false));
  }, []);

  const handleLogin = () => {
    window.location.href = `${API_BASE_URL}/auth/google/login`;
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFeedback(null);
    try {
      if (mode === 'register') {
        await registerWithEmail({
          name: form.name,
          email: form.email,
          password: form.password,
        });
        setFeedback('注册成功，已自动登录。');
      } else {
        await loginWithEmail({ email: form.email, password: form.password });
        setFeedback('登录成功。');
      }
      const current = await fetchCurrentUser();
      setUser(current);
      setForm({ ...form, password: '' });
    } catch (error) {
      setFeedback((error as Error).message || '请求失败');
    }
  };

  const googleSection = (
    <section className="auth-card">
      <h2>使用 Google 登录</h2>
      <p>连接现有 Google 账户，继续同步你的配置。</p>
      <button onClick={handleLogin}>使用 Google 登录</button>
    </section>
  );

  const emailSection = (
    <section className="auth-card">
      <div className="auth-card__header">
        <h2>{mode === 'login' ? '邮箱登录' : '邮箱注册'}</h2>
        <button
          className="link-button"
          type="button"
          onClick={() => {
            setMode(mode === 'login' ? 'register' : 'login');
            setFeedback(null);
          }}
        >
          {mode === 'login' ? '我要注册' : '我已有账号'}
        </button>
      </div>
      <form onSubmit={handleSubmit} className="auth-form">
        {mode === 'register' && (
          <label>
            昵称
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
            />
          </label>
        )}
        <label>
          邮箱
          <input
            type="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />
        </label>
        <label>
          密码
          <input
            type="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />
        </label>
        <button type="submit">{mode === 'login' ? '登录' : '注册并登录'}</button>
      </form>
      {feedback && <p className="feedback">{feedback}</p>}
    </section>
  );

  return (
    <main className="container">
      <h1>CrazyLI Login</h1>
      <p>选择你喜欢的方式登录，Google 或邮箱 + 密码。</p>
      {loading && <p>正在加载...</p>}
      {!loading && !user && (
        <div className="auth-grid">
          {googleSection}
          {emailSection}
        </div>
      )}
      {user && (
        <section>
          <p>欢迎回来！</p>
          <UserCard user={user} />
        </section>
      )}
    </main>
  );
}
