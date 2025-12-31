import { useState } from "react";
import UserForm from "./components/UserForm";
import UserList from "./components/UserList";

export default function App() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [showUsers, setShowUsers] = useState(true);

  const handleUserCreated = () => {
    setRefreshKey((k) => k + 1);
    setShowUsers(true); // auto show list after create
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 to-purple-600 p-6">
      <div className="max-w-2xl mx-auto space-y-6">
        {/* App Title */}
        <div className="bg-white rounded-2xl shadow-xl p-6 text-center">
          <h1 className="text-2xl font-bold text-gray-800">
            User Simple Management
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Create and view users easily
          </p>
        </div>
        
        {/* Create User Form */}
        <UserForm onUserCreated={handleUserCreated} />

        {/* Toggle Button */}
        <button
          onClick={() => setShowUsers((v) => !v)}
          className="text-sm text-white font-semibold underline"
        >
          {showUsers ? "Hide user list" : "Show user list"}
        </button>

        {/* User List */}
        {showUsers && <UserList refreshKey={refreshKey} />}
      </div>
    </div>
  );
}
