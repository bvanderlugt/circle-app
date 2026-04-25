import { useState, useEffect } from 'react'
import './App.css'

function CircleRow({ circle, onDelete }) {
  return (
    <li className="circle-row">
      <span className="circle-name">{circle.name}</span>
      <button
        className="remove-btn"
        onClick={() => onDelete(circle.id)}
        aria-label={`Remove ${circle.name}`}
      >
        Remove
      </button>
    </li>
  )
}

function App() {
  const [circles, setCircles] = useState([])
  const [count, setCount] = useState(0)
  const [inputName, setInputName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  async function fetchCircles() {
    try {
      const [listRes, countRes] = await Promise.all([
        fetch('/api/circles'),
        fetch('/api/circles/count'),
      ])
      const listData = await listRes.json()
      const countData = await countRes.json()
      setCircles(listData)
      setCount(countData.count)
    } catch (err) {
      console.error('Failed to fetch circles:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCircles()
  }, [])

  async function handleAdd(e) {
    e.preventDefault()
    const name = inputName.trim()
    if (!name) return

    try {
      const res = await fetch('/api/circles', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      })

      if (res.status === 201) {
        setInputName('')
        setError('')
        await fetchCircles()
      } else {
        const data = await res.json()
        setError(data.error || 'Something went wrong')
      }
    } catch (err) {
      console.error('Failed to add circle:', err)
      setError('Network error — please try again')
    }
  }

  async function handleDelete(id) {
    try {
      const res = await fetch(`/api/circles/${id}`, { method: 'DELETE' })
      if (res.status === 204) {
        await fetchCircles()
      } else {
        const data = await res.json()
        console.error('Delete failed:', data.error)
      }
    } catch (err) {
      console.error('Failed to delete circle:', err)
    }
  }

  const addDisabled = inputName.trim().length === 0

  return (
    <div className="app">
      <h1>Circles ({count})</h1>

      <form className="add-form" onSubmit={handleAdd}>
        <input
          type="text"
          className="name-input"
          value={inputName}
          onChange={(e) => setInputName(e.target.value)}
          placeholder="New circle name"
          aria-label="New circle name"
        />
        <button type="submit" className="add-btn" disabled={addDisabled}>
          Add
        </button>
      </form>

      {error && <p className="error-msg">{error}</p>}

      {loading ? (
        <p className="loading">Loading...</p>
      ) : circles.length === 0 ? (
        <p className="empty">No circles yet. Add one above.</p>
      ) : (
        <ul className="circle-list">
          {circles.map((circle) => (
            <CircleRow key={circle.id} circle={circle} onDelete={handleDelete} />
          ))}
        </ul>
      )}
    </div>
  )
}

export default App
