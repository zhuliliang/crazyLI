import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export function LoginSuccess() {
  const navigate = useNavigate();

  useEffect(() => {
    const timer = setTimeout(() => navigate('/'), 2000);
    return () => clearTimeout(timer);
  }, [navigate]);

  return (
    <main className="container">
      <h1>登录成功</h1>
      <p>已创建会话，即将跳转回首页。</p>
    </main>
  );
}
