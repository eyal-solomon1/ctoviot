import UploadVideo from "./UploadVideo.js";

function NewVideo({
  setSuccessMsg,
  setErrorMsg,
  errRef,
  successRef,
  errorMsg,
  successMsg,
  navigate,
}) {
  return (
    <>
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
    </>
  );
}

export default NewVideo;
