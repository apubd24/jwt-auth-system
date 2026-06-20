import React from 'react';
import CustomerManager from './components/CustomerManager';

function App() {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role') || 'readonly';

  return (
    <div>
      <CustomerManager token={token} role={role} />
    </div>
  );
}

export default App;
