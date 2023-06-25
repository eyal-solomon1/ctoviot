import { Link } from "react-router-dom";

const Missing = () => {
  return (
    <article>
      <h1>Oops!</h1>
      <p>Page Not Found</p>
      <div className="flexGrow">
        <Link
          style={{ color: "black", paddingTop: "1.5rem", fontSize: "1rem" }}
          to="/"
        >
          Go to Homepage
        </Link>
      </div>
    </article>
  );
};

export default Missing;
