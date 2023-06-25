import { useNavigate, Link } from "react-router-dom";
import useLogout from "../hooks/useLogout";

const Home = () => {
  const navigate = useNavigate();
  const logout = useLogout();

  const signOut = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <section>
      <h1>Hi there 👋</h1>
      <br />
      <p>
        Welcome to <i>ctoviot</i>, Your personal AI transcriber !
      </p>
      <br />
      <p>
        Using me is super simple, just upload a new video and I got you coverd
      </p>
      <br />
      <p>
        Go to <Link to="/videos">My videos</Link> and create your first AI
        powered transcribed video
      </p>
      <br />
      <Link to="/videos">
        Go to the <i>my videos</i> page
      </Link>
      <div className="flexGrow">
        <button onClick={signOut}>Sign Out</button>
      </div>
    </section>
  );
};

export default Home;
