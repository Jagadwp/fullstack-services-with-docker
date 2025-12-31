import { useEffect, useState } from "react";
import { getUsers } from "../api";

export default function UserList({ refreshKey }) {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    getUsers()
      .then(setUsers)
      .catch((err) => setError(err.message));
  }, [refreshKey]); 

  return (
    <div className="bg-white rounded-2xl shadow-xl p-6">
      <h2 className="text-xl font-semibold mb-4">Existing Users</h2>

      {error && (
        <p className="text-red-600 text-sm mb-2">{error}</p>
      )}

      <ul className="space-y-2">
        {users.map((u) => (
          <li
            key={u.id}
            className="border rounded-lg px-4 py-2 text-sm"
          >
            <p className="font-medium">{u.name}</p>
            <p className="text-gray-600">{u.email}</p>
          </li>
        ))}
      </ul>

      {users.length === 0 && (
        <p className="text-sm text-gray-500">No users yet</p>
      )}
    </div>
  );
}
