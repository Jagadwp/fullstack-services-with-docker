import { useState } from "react";
import { createUser } from "../api";

export default function UserForm({ onUserCreated }) {
  const [form, setForm] = useState({ name: "", email: "" });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const user = await createUser(form);
      onUserCreated(user);     // 🔥 notify parent
      setForm({ name: "", email: "" });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-2xl shadow-xl p-6">
      <h2 className="text-xl font-semibold mb-4">Create User</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          name="name"
          placeholder="Full Name"
          value={form.name}
          onChange={handleChange}
          className="w-full px-4 py-4 border rounded-md focus:outline-none focus:ring-2 focus:ring-brand"
          required
        />

        <input
          name="email"
          type="email"
          placeholder="Email Address"
          value={form.email}
          onChange={handleChange}
          className="w-full px-4 py-4 border rounded-md focus:outline-none focus:ring-2 focus:ring-brand"
          required
        />

        <button
          disabled={loading}
          className="w-full bg-brand text-white py-4 rounded-md font-semibold hover:opacity-90 transition"

        >
          {loading ? "Submitting..." : "Create User"}
        </button>
      </form>

      {error && (
        <p className="mt-3 text-red-600 text-sm">{error}</p>
      )}
    </div>
  );
}
