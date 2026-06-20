import React, { useEffect, useState } from 'react';
import {
  createCustomer,
  deleteCustomer,
  fetchCustomers,
  updateCustomer,
  setAuthToken,
} from '../api/customerApi';

const initialContact = { name: '', email: '', mobile: '', whatsapp: '' };
const initialPhone = { number: '', label: '' };

function CustomerManager({ token, role }) {
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [selectedId, setSelectedId] = useState(null);
  const [form, setForm] = useState({
    company_name: '',
    office_address: '',
    support_email: '',
    common_phones: [initialPhone],
    contacts: [initialContact],
  });

  useEffect(() => {
    setAuthToken(token);
  }, [token]);

  useEffect(() => {
    refreshCustomers();
  }, []);

  const refreshCustomers = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetchCustomers();
      setCustomers(response.data.customers || []);
    } catch (err) {
      setError('Unable to load customers.');
    } finally {
      setLoading(false);
    }
  };

  const resetForm = () => {
    setSelectedId(null);
    setForm({
      company_name: '',
      office_address: '',
      support_email: '',
      common_phones: [initialPhone],
      contacts: [initialContact],
    });
    setError('');
  };

  const handleInput = (field, value) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  const handlePhoneChange = (index, value) => {
    const list = [...form.common_phones];
    list[index] = { ...list[index], ...value };
    setForm((prev) => ({ ...prev, common_phones: list }));
  };

  const handleContactChange = (index, value) => {
    const list = [...form.contacts];
    list[index] = { ...list[index], ...value };
    setForm((prev) => ({ ...prev, contacts: list }));
  };

  const addPhone = () => {
    setForm((prev) => ({ ...prev, common_phones: [...prev.common_phones, initialPhone] }));
  };

  const removePhone = (index) => {
    const list = [...form.common_phones];
    list.splice(index, 1);
    setForm((prev) => ({ ...prev, common_phones: list.length ? list : [initialPhone] }));
  };

  const addContact = () => {
    setForm((prev) => ({ ...prev, contacts: [...prev.contacts, initialContact] }));
  };

  const removeContact = (index) => {
    const list = [...form.contacts];
    list.splice(index, 1);
    setForm((prev) => ({ ...prev, contacts: list.length ? list : [initialContact] }));
  };

  const handleEdit = (customer) => {
    setSelectedId(customer.id);
    setForm({
      company_name: customer.company_name || '',
      office_address: customer.office_address || '',
      support_email: customer.support_email || '',
      common_phones: customer.common_phones.length ? customer.common_phones : [initialPhone],
      contacts: customer.contacts.length ? customer.contacts : [initialContact],
    });
    setError('');
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setLoading(true);
    setError('');

    const payload = {
      company_name: form.company_name,
      office_address: form.office_address,
      support_email: form.support_email,
      common_phones: form.common_phones.filter((item) => item.number.trim()),
      contacts: form.contacts.filter((item) => item.name.trim()),
    };

    if (!payload.contacts.length) {
      setError('Please add at least one contact person.');
      setLoading(false);
      return;
    }

    try {
      if (selectedId) {
        await updateCustomer(selectedId, payload);
      } else {
        await createCustomer(payload);
      }
      await refreshCustomers();
      resetForm();
    } catch (err) {
      setError('Action failed. Please check the form and try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Delete this customer?')) {
      return;
    }
    setLoading(true);
    setError('');
    try {
      await deleteCustomer(id);
      await refreshCustomers();
      if (selectedId === id) {
        resetForm();
      }
    } catch (err) {
      setError('Failed to delete customer.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 1000, margin: '0 auto', padding: 24 }}>
      <h1>Customer Management</h1>

      {error && <div style={{ color: 'red', marginBottom: 16 }}>{error}</div>}

      <section style={{ marginBottom: 32 }}>
        <h2>{selectedId ? 'Edit Customer' : 'Add Customer'}</h2>
        <form onSubmit={handleSubmit}>
          <div>
            <label>Company Name</label>
            <input
              value={form.company_name}
              onChange={(e) => handleInput('company_name', e.target.value)}
              required
            />
          </div>
          <div>
            <label>Office Address</label>
            <input
              value={form.office_address}
              onChange={(e) => handleInput('office_address', e.target.value)}
              required
            />
          </div>
          <div>
            <label>Support Email</label>
            <input
              type="email"
              value={form.support_email}
              onChange={(e) => handleInput('support_email', e.target.value)}
              required
            />
          </div>

          <div style={{ marginTop: 16 }}>
            <h3>Common Phone Numbers</h3>
            {form.common_phones.map((phone, index) => (
              <div key={index} style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
                <input
                  placeholder="Number"
                  value={phone.number}
                  onChange={(e) => handlePhoneChange(index, { number: e.target.value })}
                  required
                />
                <input
                  placeholder="Label"
                  value={phone.label}
                  onChange={(e) => handlePhoneChange(index, { label: e.target.value })}
                />
                {form.common_phones.length > 1 && (
                  <button type="button" onClick={() => removePhone(index)}>
                    Remove
                  </button>
                )}
              </div>
            ))}
            <button type="button" onClick={addPhone}>
              Add Phone
            </button>
          </div>

          <div style={{ marginTop: 16 }}>
            <h3>Contact Persons</h3>
            {form.contacts.map((contact, index) => (
              <div key={index} style={{ border: '1px solid #ccc', padding: 12, marginBottom: 8 }}>
                <div>
                  <label>Name</label>
                  <input
                    value={contact.name}
                    onChange={(e) => handleContactChange(index, { name: e.target.value })}
                    required
                  />
                </div>
                <div>
                  <label>Email</label>
                  <input
                    type="email"
                    value={contact.email}
                    onChange={(e) => handleContactChange(index, { email: e.target.value })}
                    required
                  />
                </div>
                <div>
                  <label>Mobile</label>
                  <input
                    value={contact.mobile}
                    onChange={(e) => handleContactChange(index, { mobile: e.target.value })}
                    required
                  />
                </div>
                <div>
                  <label>WhatsApp</label>
                  <input
                    value={contact.whatsapp}
                    onChange={(e) => handleContactChange(index, { whatsapp: e.target.value })}
                  />
                </div>
                {form.contacts.length > 1 && (
                  <button type="button" onClick={() => removeContact(index)}>
                    Remove Contact
                  </button>
                )}
              </div>
            ))}
            <button type="button" onClick={addContact}>
              Add Contact Person
            </button>
          </div>

          {role === 'admin' ? (
            <button type="submit" disabled={loading} style={{ marginTop: 16 }}>
              {selectedId ? 'Update Customer' : 'Create Customer'}
            </button>
          ) : (
            <div style={{ marginTop: 16, color: '#555' }}>
              Read-only users cannot create or edit customers.
            </div>
          )}
          {selectedId && role === 'admin' && (
            <button type="button" onClick={resetForm} style={{ marginLeft: 12 }}>
              Cancel
            </button>
          )}
        </form>
      </section>

      <section>
        <h2>Customers</h2>
        {loading && <div>Loading...</div>}
        {!loading && customers.length === 0 && <div>No customers found.</div>}
        {!loading && customers.length > 0 && (
          <table width="100%" border="1" cellPadding="8" style={{ borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th>Company</th>
                <th>Address</th>
                <th>Support Email</th>
                <th>Contacts</th>
                <th>Phones</th>
                {role === 'admin' && <th>Actions</th>}
              </tr>
            </thead>
            <tbody>
              {customers.map((customer) => (
                <tr key={customer.id}>
                  <td>{customer.company_name}</td>
                  <td>{customer.office_address}</td>
                  <td>{customer.support_email}</td>
                  <td>
                    {customer.contacts.map((contact) => (
                      <div key={contact.id}>
                        <strong>{contact.name}</strong> / {contact.email} / {contact.mobile}
                        {contact.whatsapp ? ` / WhatsApp: ${contact.whatsapp}` : ''}
                      </div>
                    ))}
                  </td>
                  <td>
                    {customer.common_phones.map((phone) => (
                      <div key={phone.id}>
                        {phone.number} {phone.label ? `(${phone.label})` : ''}
                      </div>
                    ))}
                  </td>
                  {role === 'admin' && (
                    <td>
                      <button type="button" onClick={() => handleEdit(customer)}>
                        Edit
                      </button>
                      <button type="button" onClick={() => handleDelete(customer.id)} style={{ marginLeft: 8 }}>
                        Delete
                      </button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </div>
  );
}

export default CustomerManager;
