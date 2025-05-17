import { StrictMode } from 'react'
import { createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {BrowserRouter, Route, Routes} from "react-router-dom";
import LoginPage from "./pages/Login";
import RegisterPage from "./pages/Register";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <BrowserRouter>
          <Routes>
              <Route path={"/"} element={<App />}></Route>
              <Route path={"/login"} element={<LoginPage />}></Route>
              <Route path={"/register"} element={<RegisterPage />}></Route>
          </Routes>
      </BrowserRouter>
  </StrictMode>,
)
