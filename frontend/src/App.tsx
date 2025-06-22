import "./App.css";

import { useEffect, useState } from "react";

interface Todo {
  id: number;
  title: string;
  done: boolean;
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState<string>("");
  // const API_URL = import.meta.env.VITE_API_URL;

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    const res = await fetch(`/api/todos`);
    const data = await res.json();
    setTodos(data);
  };

  const addTodo = async () => {
    const res = await fetch(`/api/todos`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title: newTodo, done: false }),
    });
    if (res.ok) {
      setNewTodo("");
      fetchTodos();
    }
  };

  const toggleTodo = async (id: number) => {
    await fetch(`/api/todos/${id}/toggle`, {
      method: "POST",
    });
    fetchTodos();
  };

  const deleteTodo = async (id: number) => {
    await fetch(`/api/todos/${id}`, {
      method: "DELETE",
    });
    setTodos(todos.filter((todo) => todo.id !== id));
  };

  return (
    <div className="min-h-screen bg-gray-100 p-10">
      <div className="max-w-md mx-auto bg-white p-6 rounded shadow">
        <h1 className="text-2xl font-bold mb-4">Todo App</h1>
        <div className="flex mb-4">
          <input
            className="border p-2 flex-1 rounded-l"
            type="text"
            value={newTodo}
            onChange={(e) => setNewTodo(e.target.value)}
            placeholder="Add new task"
          />
          <button
            onClick={addTodo}
            className="bg-blue-500 text-white px-4 rounded-r"
          >
            Add
          </button>
        </div>
        <ul>
          {todos.map((todo) => (
            <li
              key={todo.id}
              className="flex justify-between items-center py-2 border-b"
            >
              <span
                className={`flex-1 ${todo.done ? "line-through text-gray-400" : ""}`}
              >
                {todo.title}
              </span>
              <button
                onClick={() => toggleTodo(todo.id)}
                className="text-sm text-blue-500"
              >
                {todo.done ? "Undo" : "Done"}
              </button>
              <button onClick={() => deleteTodo(todo.id)}>Delete</button>
            </li>
          ))}
        </ul>
      </div>
      <span>
        ⚙️ Powered by <strong>Argo CD</strong> + <strong>Jenkins</strong>
      </span>
    </div>
  );
}

export default App;
