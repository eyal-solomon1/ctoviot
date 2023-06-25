import { useState, useEffect, useRef } from "react";
import { Link } from "react-router-dom";
import useAxiosPrivate from "../hooks/useAxiosPrivate";
import VideoC from "../components/Video";
import { useNavigate, useLocation, Navigate } from "react-router-dom";
import { faFolderPlus } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const Video = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [errMsg, setErrMsg] = useState("");
  const [videos, setVideos] = useState();
  const axiosPrivate = useAxiosPrivate();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    let isMounted = true;
    const controller = new AbortController();

    const getVideos = async () => {
      try {
        const response = await axiosPrivate.get("/videos", {
          signal: controller.signal,
        });
        isMounted && setVideos(response.data?.payload?.videos);
      } catch (err) {
        navigate("/login", { state: { from: location }, replace: true });
      } finally {
        isMounted && setIsLoading(false);
      }
    };

    getVideos();

    return () => {
      isMounted = false;
      controller.abort();
    };
  }, []);

  const handleNewVideo = () => {
    navigate("/new-video");
  };

  return (
    <section>
      <header className="new-video">
        <h1>Your videos 📹 </h1>
        <div className="new-video-icon">
          <button onClick={() => handleNewVideo()}>
            <FontAwesomeIcon icon={faFolderPlus} />
          </button>

          <p>New video</p>
        </div>
      </header>
      <p className={errMsg ? "errmsg" : "offscreen"} aria-live="assertive">
        {errMsg}
      </p>{" "}
      <br />
      <div className="videos">
        {videos?.length
          ? videos.map((vid, index) => (
              <VideoC
                videoIdentifier={vid.video_identifier}
                videoDesc={vid.video_decs}
                videoName={vid.video_name}
                key={index}
                videoRemotePath={vid.video_remote_path}
                navigate={navigate}
                setErrMsg={setErrMsg}
              />
            ))
          : !isLoading && <Navigate to="/new-video" replace />}
      </div>
      <div className="flexGrow">
        <Link to="/">Home</Link>
      </div>
    </section>
  );
};

export default Video;
