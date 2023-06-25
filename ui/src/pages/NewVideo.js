import { useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { faVideo } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import UploadVideo from "../components/UploadVideo.js";

const NewVideo = () => {
  const [errorMsg, setErrorMsg] = useState("");
  const errRef = useRef();
  const [successMsg, setSuccessMsg] = useState("");
  const successRef = useRef();
  const navigate = useNavigate();

  const handleBackNavigate = () => {
    navigate(-1);
  };
  return (
    <section>
      <header className="new-video">
        <h1>Upload a video 📹 </h1>
        <div className="new-video-icon">
          <button onClick={() => handleBackNavigate()}>
            <FontAwesomeIcon icon={faVideo} />
          </button>
          <p>My videos</p>
        </div>
      </header>
      <br />
      <p
        ref={errRef}
        className={errorMsg ? "errmsg" : "offscreen"}
        aria-live="assertive"
      >
        {errorMsg}
      </p>{" "}
      <p
        ref={successRef}
        className={successMsg ? "successmsg" : "offscreen"}
        aria-live="assertive"
      >
        {successMsg}
      </p>
      <p>Unfortunately, there are no videos to display ..</p>
      <br />
      <p>Go ahead and upload your first one !</p>
      <UploadVideo
        setSuccessMsg={setSuccessMsg}
        successRef={successRef}
        setErrorMsg={setErrorMsg}
        errRef={errRef}
        navigate={navigate}
      />
    </section>
  );
};

export default NewVideo;
