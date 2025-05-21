import { StrictMode } from 'react'
import { createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Login from "./pages/Login";
import Home from "./pages/Home";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <BrowserRouter>
          <Routes>
              <Route path={"/"} element={<App />}></Route>
              <Route path={"/login"} element={<Login />}></Route>
              <Route path={"/home"} element={<Home />}></Route>
          </Routes>
      </BrowserRouter>
  </StrictMode>,
)
