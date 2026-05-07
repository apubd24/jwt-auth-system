import React, { useState, useEffect } from 'react';
import api from '../services/api';

const DeviceForm = ({ device, onSave, onCancel }) => {
  const [form, setForm] = useState({ name: '', serial: '', description: '' });

  useEffect(() => {
    if (device && device.id) {
      setForm(device);
    } else {
      setForm({ name: '', serial: '', description: '' });
    }
  }, [device]);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (form.id) {
      await api.put(`/devices/${form.id}`, form);
    } else {
      await api.post('/devices', form);
    }
    onSave();
  };

  return (
    <form onSubmit={handleSubmit} style={{ border: '1px solid #ccc', padding: '1rem', margin: '1rem 0' }}>
      <h3>{form.id ? 'Edit Device' : 'New Device'}</h3>
      <input name="name" placeholder="Name" value={form.name} onChange={handleChange} required />
      <input name="serial" placeholder="Serial" value={form.serial} onChange={handleChange} required />
      <input name="description" placeholder="Description" value={form.description} onChange={handleChange} />
      <button type="submit">Save</button>
      <button type="button" onClick={onCancel}>Cancel</button>
    </form>
  );
};

export default DeviceForm;