import Register from "./pages/Register";
import Login from "./pages/Login";
import Home from "./pages/Home";
import Layout from "./pages/Layout";
import Video from "./pages/Video";
import Missing from "./pages/Missing";
import Unauthorized from "./pages/Unauthorized";
import RequireAuth from "./pages/RequireAuth";
import Out from "./pages/Out";
import PersistLogin from "./pages/PersistLogin";
import { Routes, Route } from "react-router-dom";
import NewVideo from "./pages/NewVideo";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route path="/" element={<Out />} />
        {/* public routes */}
        {/* <Route path="login" lement={<Login />} />
        <Route path="register" element={<Register />} />
        <Route path="unauthorized" element={<Unauthorized />} /> */}

        {/* Protected routes */}
        {/* <Route element={<PersistLogin />}>
          <Route element={<RequireAuth />}>
            <Route path="/" element={<Home />} />
            <Route path="videos" element={<Video />} />
            <Route path="new-video" element={<NewVideo />} />
          </Route>
        </Route> */}

        {/* catch all */}
        <Route path="*" element={<Missing />} />
      </Route>
    </Routes>
  );
}

export default App;
