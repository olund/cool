import {useCallback, useEffect, useState} from '@lynx-js/react'

import './App.css'
import type {BaseTouchEvent, Target} from "@lynx-js/types";

interface Todo {
    id: number;
    name: string;
    description: string;
    done: boolean;
}

const serverAddr = "http://192.168.50.98:8080";

export function App() {
    const [todos, setTodos] = useState<Todo[]>([]);
    const [newTodo, setNewTodo] = useState("")

    const toggleDone = useCallback((event: BaseTouchEvent<Target>, todo: Todo) => {
        const updatedTodos = todos.map(todox => {
            if (todox.id === todo.id) {

                const newTodo = {
                    id: todo.id,
                    name: todox.name,
                    description: todox.description,
                    done: !todox.done
                }

                // updateTodo(newTodo)
                //     .then((val) => console.log("updated todo db", val))

                return newTodo
            } else {
                return todox
            }
        })

         setTodos(updatedTodos)
    }, [todos])

    const addTodo = (e: BaseTouchEvent<Target>) => {
        console.log("adding new todo!", "new todo:", newTodo)
        if (newTodo.trim() === "") return

        const newTask = {
            id: todos.length + 1, // This should be handled by the backend in a real app
            name: newTodo,
            description: "a value here",
            done: false
        }

        setTodos([...todos, newTask])
        setNewTodo("") // Clear input after adding
    }

    const getTodos = async () => {
        try {
            const json = await fetch(serverAddr + "/todo",)
                .then((res) => res.json())
                .catch(err => console.error(err));
            if (json.length > 0) {
                setTodos(json);
            }
        } catch (err) {
            console.error(err);
        }
    };

    const updateTodo = async (newTodo: Todo) => {
        try {
            const json = await fetch(serverAddr + "/todo/" + newTodo.id, {
                method: "PUT",
                body: JSON.stringify(newTodo),
                headers: {"Content-Type": "application/json"}
            })
                .then((res) => res.json())
                .catch(err => console.error(err));
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        getTodos()
    }, [])

    return (
        <view className="App">
                <text className="Title">Todos</text>

                <scroll-view
                    className="TodoList"
                    scroll-orientation="vertical"
                >
                    {!!todos && todos.map(todo => (
                        <view
                            key={`todo-list-item-${todo.id}`}
                            className="TodoItem"
                            bindtap={(e) => toggleDone(e, todo)}
                        >
                            <text className="TodoItemText">{todo.name}</text>
                            <text className={todo.done ? "TodoItemDone" : "TodoItemNotDone"}>
                                {todo.done ? "✔" : "✖"}
                            </text>
                        </view>
                    ))}
                </scroll-view>

                {/* Input Field for New Todo */}
                <view className="InputContainer">
                    <input
                        className="InputField"
                        value={newTodo}
                        onChange={(e) => {
                            console.log("onChange", e.target.value);
                            setNewTodo(e.target.value)
                        }}

                        placeholder="Enter new todo..."
                    />
                    <text className="AddButton" bindtap={(e) => addTodo(e)}>Add</text>
                </view>
        </view>
    )
}
