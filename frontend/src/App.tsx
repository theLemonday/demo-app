// ---- client/src/App.tsx ----
import "./App.css";
import { useEffect, useState } from "react";

interface Todo {
  id: number;
  title: string;
  done: boolean;
}

function App() {
  const [error, setError] = useState("");
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState<string>("");
  const [username, setUsername] = useState(
    localStorage.getItem("username") || "",
  );
  const [password, setPassword] = useState(
    localStorage.getItem("password") || "",
  );
  const [isLoggedIn, setIsLoggedIn] = useState(
    !!localStorage.getItem("username"),
  );
  const credentials = btoa(`${username}:${password}`);

  useEffect(() => {
    if (!isLoggedIn) return;
    fetchTodos().catch(() => logout());
  }, [isLoggedIn]);

  useEffect(() => {
    fetchTodos();
  }, []);

  const login = () => {
    localStorage.setItem("username", username);
    localStorage.setItem("password", password);
    setIsLoggedIn(true);
  };

  const logout = () => {
    localStorage.removeItem("username");
    localStorage.removeItem("password");
    setUsername("");
    setPassword("");
    setIsLoggedIn(false);
    setTodos([]);
  };

  const fetchTodos = async () => {
    const res = await fetch(`/api/todos`, {
      headers: {
        Authorization: `Basic ${credentials}`,
      },
    });
    const data = await res.json();
    setTodos(data);
  };

  const addTodo = async () => {
    const res = await fetch(`/api/todos`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Basic ${credentials}`,
      },
      body: JSON.stringify({ title: newTodo, done: false }),
    });

    if (res.status === 403) {
      setError("Forbidden: You don't have permission to add todos.");
      return;
    }

    if (res.ok) {
      const newTodo = await res.json();
      setTodos([...todos, newTodo]);
      setNewTodo("");
    }
  };

  const deleteTodo = async (id: number) => {
    const res = await fetch(`/api/todos/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Basic ${credentials}`,
      },
    });

    if (res.status === 403) {
      setError("Forbidden: You don't have permission to delete todos.");
      return;
    }

    setTodos(todos.filter((todo) => todo.id !== id));
  };

  const toggleTodo = async (id: number) => {
    const res = await fetch(`/api/todos/${id}/toggle`, {
      method: "POST",
      headers: {
        Authorization: `Basic ${credentials}`,
      },
    });

    if (res.status === 403) {
      setError("Forbidden: You don't have permission to modify todos.");
      return;
    }

    fetchTodos();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-100 to-white p-10">
      <dialog
        open={!isLoggedIn}
        className="backdrop:bg-black/50 rounded-xl shadow-lg"
      >
        <form
          method="dialog"
          className="bg-white rounded-lg p-6 space-y-4 w-80"
        >
          <h2 className="text-lg font-semibold text-center">🔐 Sign In</h2>
          <input
            className="border w-full p-2 rounded focus:outline-none focus:ring focus:border-blue-300"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            className="border w-full p-2 rounded focus:outline-none focus:ring focus:border-blue-300"
            placeholder="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button
            type="button"
            onClick={login}
            className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
          >
            Login
          </button>
        </form>
      </dialog>

      {isLoggedIn && (
        <div className="max-w-md mx-auto bg-white p-6 rounded-xl shadow-md">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-2xl font-bold text-gray-700">📋 Todo App</h1>
            <button
              onClick={logout}
              className="text-sm text-red-500 hover:underline"
            >
              Sign Out
            </button>
          </div>

          <div className="flex mb-4">
            <input
              className="border p-2 flex-1 rounded-l focus:outline-none focus:ring focus:border-blue-300"
              type="text"
              value={newTodo}
              onChange={(e) => setNewTodo(e.target.value)}
              placeholder="Add new task"
            />
            <button
              onClick={addTodo}
              className="bg-blue-500 text-white px-4 rounded-r hover:bg-blue-600"
            >
              Add
            </button>
          </div>

          {error && (
            <div className="bg-red-100 text-red-600 border border-red-300 px-4 py-2 rounded mb-4">
              {error}
            </div>
          )}

          <ul className="space-y-2">
            {todos.map((todo) => (
              <li
                key={todo.id}
                className="flex items-center justify-between bg-gray-50 px-4 py-2 rounded border"
              >
                <span
                  className={`flex-1 ${todo.done ? "line-through text-gray-400" : "text-gray-800"}`}
                >
                  {todo.title}
                </span>
                <div className="flex gap-2 ml-4">
                  <button
                    onClick={() => toggleTodo(todo.id)}
                    className="text-xs text-blue-500 hover:underline"
                  >
                    {todo.done ? "Undo" : "Done"}
                  </button>
                  <button
                    onClick={() => deleteTodo(todo.id)}
                    className="text-xs text-red-500 hover:underline"
                  >
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="mt-10 text-center text-sm text-gray-500">
        ⚙️ Powered by <strong>Argo CD</strong> + <strong>Jenkins</strong>
      </div>
    </div>
  );
}

export default App;
