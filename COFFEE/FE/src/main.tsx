// Harus paling awal: memindahkan key localStorage lama ke nama baru sebelum store dibuat.

import React from 'react'
import ReactDOM from 'react-dom/client'

import { App } from './app/App'
import '@/styles/globals.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
