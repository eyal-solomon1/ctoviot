import { useState, useEffect } from "react";
import useAxiosPrivate from "../hooks/useAxiosPrivate";
import { faTimes, faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const VIDEO_NAME_REGEX = /^.{4,18}$/;
const VIDEO_DESC_REGEX = /^.{8,45}$/;

const validFilesTypes = `
video/x-msvideo,
video/x-flv,
video/x-matroska,
video/x-ms-wmv,
video/x-ms-asf,
video/x-avi,
video/x-mpeg,
video/x-mpeg2a,
video/mp4,
video/x-quicktime,
video/x-sgi-movie,
video/x-m4v,
video/x-webex,
video/x-mng,
video/x-nsv,
video/x-ogg,
video/x-matroska,
video/x-f4v,
video/x-fli,
video/x-vob,
video/quicktime
`;

const Video = ({
  successRef,
  setSuccessMsg,
  setErrorMsg,
  errRef,
  navigate,
}) => {
  const [selectedFile, setSelectedFile] = useState(null);

  const [videoName, setVideoName] = useState("");
  const [validVideoName, setValidVideoName] = useState(false);
  const [videoNameFocus, setVideoNameFocus] = useState(false);

  const [description, setDescription] = useState("");
  const [validDescription, setValidDescription] = useState(false);
  const [descriptionFocus, setDescriptionFocus] = useState(false);

  const [loading, setLoading] = useState(false);

  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    setValidVideoName(VIDEO_NAME_REGEX.test(videoName));
  }, [videoName]);

  useEffect(() => {
    setValidDescription(VIDEO_DESC_REGEX.test(description));
  }, [description]);

  const handleRemoveFile = () => {
    setSelectedFile(null);
  };

  const handleShortNameChange = (event) => {
    setVideoName(event.target.value);
  };

  const handleDescriptionChange = (event) => {
    setDescription(event.target.value);
  };

  const onChangeHandler = async (event) => {
    const file = event.target.files[0];

    const video = document.createElement("video");
    video.preload = "metadata";

    video.onloadedmetadata = () => {
      if (video.duration > 15) {
        setErrorMsg("video length is can't be more then 15s");
        return;
      }
      setSelectedFile(file);
      setErrorMsg("");
    };

    video.src = URL.createObjectURL(file);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData();
    formData.append("inputFile", selectedFile);
    formData.append("videoName", videoName);
    formData.append("description", description);

    try {
      setLoading(true);
      await axiosPrivate.post("/new_video", formData);
      setSuccessMsg("Video upload successfuly ! refreshing page ..");
      successRef.current.focus();
      setErrorMsg("");
      setLoading(false);

      await new Promise((res) => setTimeout(res, 2500));
      navigate("/videos");
    } catch (err) {
      const returnedErrorMsg =
        err.response?.data?.error?.error || "Uplaod failed";
      console.error("here", returnedErrorMsg);

      if (!err?.response) {
        setErrorMsg("No Server Response");
      } else if (err.response?.status === 401) {
        setErrorMsg("Unauthorized");
      } else {
        setErrorMsg(returnedErrorMsg);
      }
      errRef.current.focus();
      setSuccessMsg("");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="upload-container">
      <form onSubmit={handleSubmit} className="upload-form">
        <label htmlFor="file-upload" className="custom-file-upload">
          <input
            type="file"
            accept={validFilesTypes}
            id="file-upload"
            name="file"
            onChange={onChangeHandler}
            className="file-input"
          />
          Choose File
        </label>
        {selectedFile && (
          <div className="selected-file">
            <p className="file-name">{selectedFile.name}</p>
            <FontAwesomeIcon onClick={handleRemoveFile} icon={faTimes} />
          </div>
        )}
        <input
          placeholder="Funny Video name"
          value={videoName}
          onChange={handleShortNameChange}
          className="upload-input"
          onFocus={() => setVideoNameFocus(true)}
          onBlur={() => setVideoNameFocus(false)}
        />
        <p
          id="namenote"
          className={
            videoNameFocus && !validVideoName ? "instructions" : "offscreen"
          }
        >
          <FontAwesomeIcon icon={faInfoCircle} />
          At least 4 characters.
          <br />
          Up to 18 characters.
        </p>
        <textarea
          placeholder="A short description"
          value={description}
          onChange={handleDescriptionChange}
          className="upload-textarea"
          onFocus={() => setDescriptionFocus(true)}
          onBlur={() => setDescriptionFocus(false)}
        ></textarea>
        <p
          id="descnote"
          className={
            descriptionFocus && !validDescription ? "instructions" : "offscreen"
          }
        >
          <FontAwesomeIcon icon={faInfoCircle} />
          At least 8 characters.
          <br />
          Up to 45 characters.
        </p>
        {loading && <span className="loading">Loading...</span>}
        <button
          disabled={
            !validDescription || !validVideoName || !selectedFile || loading
          }
          type="submit"
          className="upload-button"
        >
          Upload
        </button>
      </form>
    </div>
  );
};

export default Video;
