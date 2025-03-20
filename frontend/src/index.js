import React, { useState } from 'react'; import { createRoot } from 'react-dom/client'; import axios from 'axios';

const App = () => {
  const [form, setForm] = useState({ nome: '', email: '', cpf: '', nascimento: '', senha: '' });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const register = async () => {
    try {
      await axios.post('http://localhost:8080/register', form);
      alert('Usuário cadastrado com sucesso!');
    } catch (error) {
      alert('Erro ao cadastrar usuário');
    }
  };

  return (
    <div className='p-4 max-w-md mx-auto'>
      <h1 className='text-xl font-bold mb-4'>Cadastro de Usuário</h1>
      <input type='text' name='nome' placeholder='Nome' className='border p-2 w-full mb-2' onChange={handleChange} />
      <input type='email' name='email' placeholder='Email' className='border p-2 w-full mb-2' onChange={handleChange} />
      <input type='text' name='cpf' placeholder='CPF' className='border p-2 w-full mb-2' onChange={handleChange} />
      <input type='date' name='nascimento' className='border p-2 w-full mb-2' onChange={handleChange} />
      <input type='password' name='senha' placeholder='Senha' className='border p-2 w-full mb-2' onChange={handleChange} />
      <button className='bg-blue-500 text-white p-2 w-full' onClick={register}>Cadastrar</button>
    </div>
  );
};

createRoot(document.getElementById('root')).render(<App />);