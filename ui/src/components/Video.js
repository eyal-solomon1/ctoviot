import { useState } from "react";
import VideoBlock from "./VideoBlock.js";
import { faTrash } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import useAxiosPrivate from "../hooks/useAxiosPrivate.js";
import { REMOTE_FILES_ENDPOINT_PREFIX } from "../util/env.js";

function Video({
  videoName,
  videoDesc,
  videoRemotePath,
  videoIdentifier,
  navigate,
  setErrMsg,
}) {
  const [isLoading, setIsLoading] = useState(false);
  const axiosPrivate = useAxiosPrivate();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const handleDeleteVideo = async (videoName) => {
    try {
      setIsLoading(true);
      await axiosPrivate.post(
        "/delete-video",
        { video_identifier: videoIdentifier },
        { withCredentials: true }
      );
      navigate("/");
    } catch (err) {
      console.log(err);
      setErrMsg(`couldn't delete ${videoName} ..`);
    } finally {
      setIsLoading(false);
    }

    setShowDeleteDialog(false);
  };

  return (
    <div className="video">
      <div className="trash-icon">
        <FontAwesomeIcon
          icon={faTrash}
          onClick={() => setShowDeleteDialog(true)}
        />
      </div>
      <p id="name">
        <u>🏷️</u> {videoName}
      </p>
      <p id="desc">
        {" "}
        <u>
          <i>📙</i>
        </u>{" "}
        {videoDesc}
      </p>
      <VideoBlock
        videoPath={`${REMOTE_FILES_ENDPOINT_PREFIX}/${videoRemotePath}`}
      />

      {showDeleteDialog ? (
        isLoading ? (
          <span className="loading">Deleting ...</span>
        ) : (
          <div className="delete-dialog">
            <div>
              <p>Are you sure you want to delete "{videoName}" video?</p>
            </div>
            <div className="delete-dialog-btns">
              <button onClick={() => handleDeleteVideo(videoName)}>
                Delete
              </button>
              <button onClick={() => setShowDeleteDialog(false)}>Cancel</button>
            </div>
          </div>
        )
      ) : null}
    </div>
  );
}

export default Video;
