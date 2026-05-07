import React, { useEffect, useState } from 'react';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

const AdminUsers = () => {
  const [users, setUsers] = useState([]);
  const [editingUser, setEditingUser] = useState(null);
  const [newUser, setNewUser] = useState({ username: '', password: '', role: 'readonly' });
  const { user: currentUser } = useAuth();

  const fetchUsers = async () => {
    const res = await api.get('/users');
    setUsers(res.data.users);
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleCreate = async (e) => {
    e.preventDefault();
    await api.post('/users', newUser);
    setNewUser({ username: '', password: '', role: 'readonly' });
    fetchUsers();
  };

  const handleUpdate = async (id, updatedData) => {
    await api.put(`/users/${id}`, updatedData);
    fetchUsers();
    setEditingUser(null);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Delete user?')) {
      await api.delete(`/users/${id}`);
      fetchUsers();
    }
  };

  const toggleActive = (user) => {
    handleUpdate(user.id, { is_active: !user.is_active });
  };

  return (
    <div>
      <h2>Manage Users (Admin only)</h2>
      <form onSubmit={handleCreate}>
        <input placeholder="Username" value={newUser.username} onChange={e => setNewUser({...newUser, username: e.target.value})} required />
        <input type="password" placeholder="Password" value={newUser.password} onChange={e => setNewUser({...newUser, password: e.target.value})} required />
        <select value={newUser.role} onChange={e => setNewUser({...newUser, role: e.target.value})}>
          <option value="readonly">Readonly</option>
          <option value="admin">Admin</option>
        </select>
        <button type="submit">Create User</button>
      </form>
      <table border="1" cellPadding="8">
        <thead><tr><th>ID</th><th>Username</th><th>Role</th><th>Active</th><th>Actions</th></tr></thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id}>
              <td>{u.id}</td>
              <td>{u.username}</td>
              <td>{u.role}</td>
              <td>{u.is_active ? 'Yes' : 'No'}</td>
              <td>
                {editingUser?.id === u.id ? (
                  <>
                    <select value={editingUser.role} onChange={e => setEditingUser({...editingUser, role: e.target.value})}>
                      <option value="readonly">Readonly</option>
                      <option value="admin">Admin</option>
                    </select>
                    <button onClick={() => handleUpdate(u.id, { role: editingUser.role })}>Save</button>
                    <button onClick={() => setEditingUser(null)}>Cancel</button>
                  </>
                ) : (
                  <>
                    <button onClick={() => setEditingUser(u)}>Edit Role</button>
                    <button onClick={() => toggleActive(u)}>{u.is_active ? 'Deactivate' : 'Activate'}</button>
                    <button onClick={() => handleDelete(u.id)}>Delete</button>
                  </>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminUsers;