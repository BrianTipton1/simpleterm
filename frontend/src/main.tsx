import "./index.css";
import { Route, Routes } from "react-router";
import { createRoot } from "react-dom/client";
import { HashRouter } from "react-router";
import App from "./App";

const root = document.getElementById("root")!;

createRoot(root).render(
  <HashRouter basename={"/"}>
    <Routes>
      <Route path="/" element={<App />} />
    </Routes>
  </HashRouter>,
);
