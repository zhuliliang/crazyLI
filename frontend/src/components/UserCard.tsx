import type { User } from '../lib/api';

interface Props {
  user: User;
}

export function UserCard({ user }: Props) {
  return (
    <div className="user-card">
      {user.picture && (
        <img src={user.picture} alt={user.name} width={64} height={64} />
      )}
      <div>
        <p><strong>{user.name}</strong></p>
        <p>{user.email}</p>
      </div>
    </div>
  );
}
